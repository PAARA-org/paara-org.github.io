/*
 * Sun/moon dark-mode toggle. Light is the default; a stored preference is
 * applied before paint by an inline script in head-custom.html to prevent a
 * flash. This file builds the toggle button and remembers the choice.
 */
(function () {
  "use strict";

  var STORAGE_KEY = "paara-theme";

  function currentTheme() {
    return document.documentElement.getAttribute("data-theme") === "dark"
      ? "dark"
      : "light";
  }

  function applyTheme(theme) {
    if (theme === "dark") {
      document.documentElement.setAttribute("data-theme", "dark");
    } else {
      document.documentElement.removeAttribute("data-theme");
    }
    try {
      localStorage.setItem(STORAGE_KEY, theme);
    } catch (e) {
      /* localStorage may be unavailable */
    }
    var btn = document.querySelector(".theme-toggle");
    if (btn) {
      btn.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
    }
  }

  function buildButton() {
    var btn = document.createElement("button");
    btn.className = "theme-toggle";
    btn.type = "button";
    btn.setAttribute("aria-label", "Toggle dark mode");
    btn.setAttribute("aria-pressed", currentTheme() === "dark" ? "true" : "false");
    btn.title = "Toggle light / dark mode";

    // Moon shows in light mode (click to go dark); sun shows in dark mode.
    var moon =
      '<svg class="icon-moon" viewBox="0 0 24 24" aria-hidden="true">' +
      '<path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z"/></svg>';

    var sun =
      '<svg class="icon-sun" viewBox="0 0 24 24" aria-hidden="true">' +
      '<path d="M12 17a5 5 0 1 1 0-10 5 5 0 0 1 0 10zm0-13a1 1 0 0 1 1 1v1.5a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zm0 13.5a1 1 0 0 1 1 1V20a1 1 0 1 1-2 0v-1.5a1 1 0 0 1 1-1zM4 12a1 1 0 0 1 1-1h1.5a1 1 0 1 1 0 2H5a1 1 0 0 1-1-1zm13.5 0a1 1 0 0 1 1-1H20a1 1 0 1 1 0 2h-1.5a1 1 0 0 1-1-1zM6.3 6.3a1 1 0 0 1 1.4 0l1.1 1.1a1 1 0 1 1-1.4 1.4L6.3 7.7a1 1 0 0 1 0-1.4zm9 9a1 1 0 0 1 1.4 0l1.1 1.1a1 1 0 0 1-1.4 1.4l-1.1-1.1a1 1 0 0 1 0-1.4zm2.5-9a1 1 0 0 1 0 1.4l-1.1 1.1a1 1 0 1 1-1.4-1.4l1.1-1.1a1 1 0 0 1 1.4 0zm-9 9a1 1 0 0 1 0 1.4l-1.1 1.1a1 1 0 0 1-1.4-1.4l1.1-1.1a1 1 0 0 1 1.4 0z"/></svg>';

    btn.innerHTML = moon + sun;
    btn.addEventListener("click", function () {
      applyTheme(currentTheme() === "dark" ? "light" : "dark");
    });
    return btn;
  }

  function init() {
    if (document.querySelector(".theme-toggle")) {
      return;
    }
    document.body.appendChild(buildButton());
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
