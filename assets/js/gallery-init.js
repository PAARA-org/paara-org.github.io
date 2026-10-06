// assets/js/gallery-init.js
// Initializes GLightbox and provides gallery toggle functionality.
(function () {
  var lightbox = null;
  var exifDateCache = {};

  function parseExifDateFromBuffer(buf) {
    if (!buf) return null;
    var bytes = new Uint8Array(buf.slice(0, 65536));
    var str = "";
    for (var i = 0; i < bytes.length; i++) {
      var b = bytes[i];
      if (b >= 32 && b <= 126) {
        str += String.fromCharCode(b);
      } else {
        str += " ";
      }
    }
    // Match EXIF date format: YYYY:MM:DD HH:MM:SS
    var m = str.match(/\b(19\d\d|20\d\d):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2})\b/);
    if (!m) return null;

    var year = parseInt(m[1], 10);
    var month = parseInt(m[2], 10) - 1;
    var day = parseInt(m[3], 10);
    var hour = parseInt(m[4], 10);
    var min = parseInt(m[5], 10);
    var sec = parseInt(m[6], 10);
    var d = new Date(year, month, day, hour, min, sec);
    if (isNaN(d.getTime())) return null;

    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric"
    }) + " at " + d.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit"
    });
  }

  function parseFilenameDate(filename) {
    if (!filename) return null;
    var m = filename.match(/(?:^|[^0-9])(20\d\d)(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})(?:[^0-9]|$)/);
    if (!m) return null;
    var d = new Date(parseInt(m[1], 10), parseInt(m[2], 10) - 1, parseInt(m[3], 10), parseInt(m[4], 10), parseInt(m[5], 10), parseInt(m[6], 10));
    if (isNaN(d.getTime())) return null;
    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric"
    }) + " at " + d.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit"
    });
  }

  function parseExifDateFromBytes(bytes) {
    if (!bytes || bytes.length === 0) return null;
    var maxLen = Math.min(bytes.length, 65536);
    var str = "";
    for (var i = 0; i < maxLen; i++) {
      var b = bytes[i];
      if (b >= 32 && b <= 126) {
        str += String.fromCharCode(b);
      } else {
        str += " ";
      }
    }
    var m = str.match(/\b(19\d\d|20\d\d):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2})\b/);
    if (!m) return null;

    var year = parseInt(m[1], 10);
    var month = parseInt(m[2], 10) - 1;
    var day = parseInt(m[3], 10);
    var hour = parseInt(m[4], 10);
    var min = parseInt(m[5], 10);
    var sec = parseInt(m[6], 10);
    var d = new Date(year, month, day, hour, min, sec);
    if (isNaN(d.getTime())) return null;

    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric"
    }) + " at " + d.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit"
    });
  }

  function fetchAndApplyExifDate(slide, slideIndex) {
    if (!slide) return;

    var titleEl = slide.querySelector(".gslide-title");
    var descEl = slide.querySelector(".gslide-desc");
    var innerEl = slide.querySelector(".gdesc-inner");
    var imgEl = slide.querySelector(".gslide-media img");

    var title = titleEl ? titleEl.textContent.trim() : "";
    var url = imgEl ? imgEl.src : "";

    if ((!title || !url) && lightbox && lightbox.elements && typeof slideIndex === "number") {
      var elData = lightbox.elements[slideIndex];
      if (elData) {
        if (!title) title = elData.title || "";
        if (!url) url = elData.href || "";
      }
    }

    // Ensure .gslide-desc exists
    if (!descEl && innerEl) {
      descEl = document.createElement("div");
      descEl.className = "gslide-desc";
      innerEl.appendChild(descEl);
    }
    if (!descEl) return;

    function applyDate(dateStr) {
      if (!dateStr) return;
      if (descEl.querySelector(".gallery-slide-time")) return;
      if (descEl.textContent.indexOf(dateStr) !== -1 || descEl.textContent.indexOf(" at ") !== -1) return;
      var span = document.createElement("span");
      span.className = "gallery-slide-time";
      span.textContent = (descEl.textContent.trim() ? " • " : "") + dateStr;
      descEl.appendChild(span);
    }

    // 0. Check data-date from pre-rendered HTML if present
    if (lightbox && lightbox.elements && typeof slideIndex === "number") {
      var elData = lightbox.elements[slideIndex];
      if (elData && elData.node && elData.node.dataset && elData.node.dataset.date) {
        applyDate(elData.node.dataset.date);
        return;
      }
    }

    // 1. Try filename date first (instant)
    var fromFilename = parseFilenameDate(title);
    if (fromFilename) {
      applyDate(fromFilename);
      return;
    }

    if (!url || !url.startsWith("http")) return;

    // 2. Check memory cache
    if (exifDateCache[url] !== undefined) {
      if (exifDateCache[url]) {
        applyDate(exifDateCache[url]);
      }
      return;
    }

    // 3. Simple CORS stream reader (reads up to ~65KB without triggering preflight)
    fetch(url, { referrerPolicy: "no-referrer" })
      .then(function (res) {
        if (!res.ok || !res.body) return null;
        var reader = res.body.getReader();
        var chunks = [];
        var totalBytes = 0;
        function readNext() {
          return reader.read().then(function (result) {
            if (result.value) {
              chunks.push(result.value);
              totalBytes += result.value.length;
            }
            if (result.done || totalBytes >= 65536) {
              reader.cancel();
              var combined = new Uint8Array(totalBytes);
              var offset = 0;
              for (var i = 0; i < chunks.length; i++) {
                combined.set(chunks[i], offset);
                offset += chunks[i].length;
              }
              return combined;
            }
            return readNext();
          });
        }
        return readNext();
      })
      .then(function (bytes) {
        if (!bytes) {
          exifDateCache[url] = null;
          return;
        }
        var dateStr = parseExifDateFromBytes(bytes);
        exifDateCache[url] = dateStr;
        if (dateStr) {
          applyDate(dateStr);
        }
      })
      .catch(function () {
        exifDateCache[url] = null;
      });
  }

  // Displays a fresh random continuous slice of <x> items on each page load.
  // The underlying DOM order is kept 100% strictly chronological (oldest to newest, 0 to N-1).
  // Slide 1 in the album is ALWAYS the oldest photo.
  function setupRandomContinuousGalleries() {
    var grids = document.querySelectorAll(".gallery-grid");
    grids.forEach(function (grid) {
      var items = Array.prototype.slice.call(grid.querySelectorAll(".gallery-item"));
      var total = items.length;
      var limit = parseInt(grid.dataset.previewLimit, 10) || 8;

      if (total <= limit) {
        items.forEach(function (el) {
          el.classList.remove("gallery-item-extra");
          el.style.display = "block";
        });
        return;
      }

      // Pick a random starting index between 0 and (total - limit)
      // so the preview set is a continuous chronological slice of the event.
      var maxStart = total - limit;
      var startIndex = Math.floor(Math.random() * (maxStart + 1));

      // Record active slice bounds on the grid for clean expand/collapse
      grid.dataset.sliceStart = startIndex;
      grid.dataset.sliceEnd = startIndex + limit;

      // Keep original DOM order intact so GLightbox slide 1 is always the oldest photo!
      items.forEach(function (el, index) {
        if (index >= startIndex && index < startIndex + limit) {
          el.classList.remove("gallery-item-extra");
          el.style.display = "block";
        } else {
          el.classList.add("gallery-item-extra");
          el.style.display = "none";
        }
      });
    });
  }

  function initGallery() {
    if (typeof GLightbox === "undefined") {
      return;
    }

    setupRandomContinuousGalleries();

    lightbox = GLightbox({
      selector: ".glightbox",
      touchNavigation: true,
      loop: true,
      autoplayVideos: true,
      zoomable: true,
      draggable: true,
      openEffect: "zoom",
      closeEffect: "fade",
      slideEffect: "slide",
      width: "960px",
      height: "540px",
      beforeSlideLoad: function (slideObj) {
        var slide = slideObj.slide;
        var container = slide ? slide.querySelector(".ginner-container") : null;
        if (container) {
          var isVideo = slide.querySelector(".gslide-external") || (lightbox.elements && lightbox.elements[slideObj.index] && lightbox.elements[slideObj.index].type === "external");
          if (isVideo) {
            container.classList.add("has-external-video");
          }
        }
      },
      afterSlideLoad: function (slideObj) {
        var slide = slideObj.slide;
        var index = slideObj.index;
        var img = slide ? slide.querySelector("img") : null;
        if (img) {
          img.setAttribute("referrerpolicy", "no-referrer");
        }
        var iframe = slide ? slide.querySelector("iframe") : null;
        if (iframe) {
          iframe.setAttribute(
            "allow",
            "autoplay; encrypted-media; fullscreen; picture-in-picture"
          );
          iframe.setAttribute("allowfullscreen", "true");
        }
        var container = slide ? slide.querySelector(".ginner-container") : null;
        if (container && iframe) {
          container.classList.add("has-external-video");
        }

        // Display image EXIF time or filename timestamp
        fetchAndApplyExifDate(slide, index);
      },
      beforeSlideChange: function (prev, current) {
        if (prev && prev.slide) {
          var prevIframe = prev.slide.querySelector("iframe");
          if (prevIframe) {
            var src = prevIframe.src;
            prevIframe.src = "";
            prevIframe.src = src;
          }
        }
      },
      onClose: function () {
        var iframes = document.querySelectorAll(".glightbox-container iframe");
        iframes.forEach(function (iframe) {
          var src = iframe.src;
          iframe.src = "";
          iframe.src = src;
        });
      }
    });

    lightbox.on("slide_after_load", function (data) {
      if (data && data.slide) {
        fetchAndApplyExifDate(data.slide, data.index);
      }
    });

    lightbox.on("slide_changed", function (data) {
      if (data && data.current && data.current.slide) {
        fetchAndApplyExifDate(data.current.slide, data.current.index);
      }
    });

    window.openGalleryFromStart = function (folderId) {
      var container = document.querySelector(
        '.gallery-container[data-gallery-id="' + CSS.escape(folderId) + '"]'
      );
      if (!container) return;
      var firstItem = container.querySelector(".gallery-item");
      if (!firstItem) return;

      if (lightbox && lightbox.elements) {
        for (var i = 0; i < lightbox.elements.length; i++) {
          if (lightbox.elements[i].node === firstItem) {
            lightbox.openAt(i);
            return;
          }
        }
      }
      firstItem.click();
    };

    window.toggleGalleryExpand = function (folderId, btn) {
      var container = document.querySelector(
        '.gallery-container[data-gallery-id="' + CSS.escape(folderId) + '"]'
      );
      if (!container) return;

      var grid = container.querySelector(".gallery-grid");
      var items = container.querySelectorAll(".gallery-item");
      var isExpanded = container.classList.contains("is-expanded");
      var limit = grid ? (parseInt(grid.dataset.previewLimit, 10) || 8) : 8;
      var sliceStart = grid && grid.dataset.sliceStart !== undefined ? parseInt(grid.dataset.sliceStart, 10) : 0;
      var sliceEnd = grid && grid.dataset.sliceEnd !== undefined ? parseInt(grid.dataset.sliceEnd, 10) : limit;

      if (isExpanded) {
        items.forEach(function (el, index) {
          if (index >= sliceStart && index < sliceEnd) {
            el.style.display = "block";
            el.classList.remove("gallery-item-extra");
          } else {
            el.style.display = "none";
            el.classList.add("gallery-item-extra");
          }
        });
        container.classList.remove("is-expanded");
        btn.innerHTML = "View all " + btn.dataset.total + " &rarr;";
      } else {
        items.forEach(function (el) {
          el.style.display = "block";
          el.classList.remove("gallery-item-extra");
        });
        container.classList.add("is-expanded");
        btn.innerHTML = "Show less &uarr;";

        // Reload lightbox so it recognizes all newly visible slides if needed
        if (lightbox && typeof lightbox.reload === "function") {
          lightbox.reload();
        }
      }
    };
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initGallery);
  } else {
    initGallery();
  }
})();
