// script/sync_photos_test.go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Media Type Classification Tests
// ============================================================================

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"photo.jpg", true},
		{"photo.JPEG", true},
		{"IMAGE.PNG", true},
		{"banner.webp", true},
		{"photo.heic", true},
		{"photo.HEIC", true},
		{"bitmap.bmp", true},
		{"photo.tif", true},
		{"photo.tiff", true},
		{"video.mp4", false},
		{"video.mov", false},
		{"document.pdf", false},
		{"archive.zip", false},
		{"readme.md", false},
		{"no_extension", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := isImageFile(tt.filename); got != tt.expected {
				t.Errorf("isImageFile(%q) = %v; want %v", tt.filename, got, tt.expected)
			}
		})
	}
}

func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"clip.mp4", true},
		{"movie.MOV", true},
		{"timelapse.webm", true},
		{"recording.m4v", true},
		{"video.mkv", true},
		{"video.avi", true},
		{"photo.jpg", false},
		{"photo.heic", false},
		{"document.pdf", false},
		{"audio.mp3", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := isVideoFile(tt.filename); got != tt.expected {
				t.Errorf("isVideoFile(%q) = %v; want %v", tt.filename, got, tt.expected)
			}
		})
	}
}

// ============================================================================
// Date & EXIF Parsing Tests
// ============================================================================

func TestParseFilenameDate(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantISO   string
		wantHuman string
	}{
		{
			name:      "standard timestamp format",
			filename:  "20261003_142311.jpg",
			wantISO:   "2026-10-03T14:23:11",
			wantHuman: "Oct 3, 2026 at 2:23 PM",
		},
		{
			name:      "timestamp with prefix",
			filename:  "IMG_20250524_113544.jpg",
			wantISO:   "2025-05-24T11:35:44",
			wantHuman: "May 24, 2025 at 11:35 AM",
		},
		{
			name:      "timestamp with suffix",
			filename:  "20260418_091500_HDR.heic",
			wantISO:   "2026-04-18T09:15:00",
			wantHuman: "Apr 18, 2026 at 9:15 AM",
		},
		{
			name:      "iPhone standard camera name",
			filename:  "IMG_7756.HEIC",
			wantISO:   "",
			wantHuman: "",
		},
		{
			name:      "Canon camera name",
			filename:  "CAL02363.jpg",
			wantISO:   "",
			wantHuman: "",
		},
		{
			name:      "invalid date month 13",
			filename:  "20261303_142311.jpg",
			wantISO:   "",
			wantHuman: "",
		},
		{
			name:      "invalid date hour 25",
			filename:  "20261003_252311.jpg",
			wantISO:   "",
			wantHuman: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iso, human := parseFilenameDate(tt.filename)
			if iso != tt.wantISO || human != tt.wantHuman {
				t.Errorf("parseFilenameDate(%q) = (%q, %q); want (%q, %q)",
					tt.filename, iso, human, tt.wantISO, tt.wantHuman)
			}
		})
	}
}

func TestParseExifDate(t *testing.T) {
	tests := []struct {
		name      string
		buf       []byte
		wantISO   string
		wantHuman string
	}{
		{
			name:      "exact EXIF date string",
			buf:       []byte("2026:10:03 08:24:21"),
			wantISO:   "2026-10-03T08:24:21",
			wantHuman: "Oct 3, 2026 at 8:24 AM",
		},
		{
			name: "EXIF date embedded in binary data",
			buf: append(
				[]byte{0xFF, 0xD8, 0xFF, 0xE1, 0x00, 0x18, 'E', 'x', 'i', 'f', 0x00, 0x00},
				[]byte("...DateTimeOriginal...2026:10:03 15:45:00...")...,
			),
			wantISO:   "2026-10-03T15:45:00",
			wantHuman: "Oct 3, 2026 at 3:45 PM",
		},
		{
			name:      "empty buffer",
			buf:       []byte{},
			wantISO:   "",
			wantHuman: "",
		},
		{
			name:      "buffer with no EXIF date",
			buf:       []byte("Exif header without any timestamp tags"),
			wantISO:   "",
			wantHuman: "",
		},
		{
			name:      "buffer with invalid date",
			buf:       []byte("2026:02:30 12:00:00"),
			wantISO:   "",
			wantHuman: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iso, human := parseExifDate(tt.buf)
			if iso != tt.wantISO || human != tt.wantHuman {
				t.Errorf("parseExifDate() = (%q, %q); want (%q, %q)",
					iso, human, tt.wantISO, tt.wantHuman)
			}
		})
	}
}

// ============================================================================
// Sorting Tests
// ============================================================================

func TestSortGalleryItems(t *testing.T) {
	items := []GalleryItem{
		{ID: "3", Name: "CAL02363.jpg", Timestamp: "2026-10-03T09:16:44"},
		{ID: "1", Name: "IMG_7756.HEIC", Timestamp: "2026-10-03T08:24:21"},
		{ID: "5", Name: "zebra.jpg", Timestamp: ""},
		{ID: "2", Name: "IMG_7757.HEIC", Timestamp: "2026-10-03T09:22:36"},
		{ID: "4", Name: "alpha.jpg", Timestamp: ""},
	}

	sortGalleryItems(items)

	expectedOrder := []string{"1", "3", "2", "4", "5"}
	for i, it := range items {
		if it.ID != expectedOrder[i] {
			t.Errorf("Index %d: got ID %q, want ID %q (item name: %s, timestamp: %s)",
				i, it.ID, expectedOrder[i], it.Name, it.Timestamp)
		}
	}
}

// ============================================================================
// HTML Scraper Parsing Tests
// ============================================================================

func TestParseEmbeddedFolderViewHTML(t *testing.T) {
	t.Run("valid entries and duplicates", func(t *testing.T) {
		sampleHTML := `
		<!DOCTYPE html>
		<html>
		<body>
			<div class="flip-entry" id="entry-1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M">
				<div class="flip-entry-title">IMG_7756.HEIC</div>
			</div>
			<div class="flip-entry" id="entry-1X7B8sa8REHGBw8mC9BH9r5ne6pVJisfG">
				<div class="flip-entry-title">timelapse.mp4</div>
			</div>
			<div class="flip-entry" id="entry-1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M">
				<div class="flip-entry-title">IMG_7756.HEIC</div>
			</div>
		</body>
		</html>
		`

		items := parseEmbeddedFolderViewHTML(sampleHTML)
		if len(items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(items))
		}

		var imgItem, vidItem *GalleryItem
		for i := range items {
			if items[i].Type == "image" {
				imgItem = &items[i]
			} else if items[i].Type == "video" {
				vidItem = &items[i]
			}
		}

		if imgItem == nil || imgItem.ID != "1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M" {
			t.Errorf("Expected image item not found or incorrect: %+v", imgItem)
		}
		if vidItem == nil || vidItem.ID != "1X7B8sa8REHGBw8mC9BH9r5ne6pVJisfG" {
			t.Errorf("Expected video item not found or incorrect: %+v", vidItem)
		}
	})

	t.Run("empty HTML returns nil", func(t *testing.T) {
		if items := parseEmbeddedFolderViewHTML("<div>No entries here</div>"); items != nil {
			t.Errorf("Expected nil items for empty HTML, got %+v", items)
		}
	})
}

func TestParsePublicScraperHTML(t *testing.T) {
	t.Run("valid aria-labels and session suffixes", func(t *testing.T) {
		sampleHTML := `
		<div aria-label="IMG_7756.HEIC Image" ssk='1:2:1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M-12-34'></div>
		<div aria-label="field_day.mp4 Video" ssk='1:2:1X7B8sa8REHGBw8mC9BH9r5ne6pVJisfG-56-78'></div>
		<div aria-label="IMG_7756.HEIC Image" ssk='1:2:1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M-99-99'></div>
		`

		items := parsePublicScraperHTML(sampleHTML)
		if len(items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(items))
		}

		if items[0].ID != "1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M" && items[1].ID != "1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M" {
			t.Errorf("Clean ID session suffix not stripped properly: %+v", items)
		}
	})

	t.Run("empty HTML returns nil", func(t *testing.T) {
		if items := parsePublicScraperHTML("<html></html>"); items != nil {
			t.Errorf("Expected nil items for empty HTML, got %+v", items)
		}
	})
}

// ============================================================================
// Markdown Target Discovery Tests
// ============================================================================

func TestDiscoverTargetsFromReader(t *testing.T) {
	mdContent := `
# Photos

## PAARA In The Park October 2026
{% include gallery.html folder_id="1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82" limit=8 %}

## Electronics Flea Market September 2026
Click [here](https://drive.google.com/drive/folders/1cRbEuwo7-tZYFPgPamXnm-8kGspkL6Z1?usp=drive_link) for more pictures.

## Field Day 2019 Timelapse
Click [here](https://drive.google.com/file/d/1X7B8sa8REHGBw8mC9BH9r5ne6pVJisfG/view?usp=drive_link) to view video.
`

	existing := map[string]Gallery{
		"1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82": {
			Title:    "Existing October 2026",
			FolderID: "1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82",
		},
		"previously-indexed-folder": {
			Title:      "Previously Indexed",
			FolderID:   "previously-indexed-folder",
			DriveURL:   "https://drive.google.com/drive/folders/previously-indexed-folder",
			TotalCount: 5,
		},
	}

	t.Run("discover all targets including existing map", func(t *testing.T) {
		r := strings.NewReader(mdContent)
		targets := discoverTargetsFromReader(r, "", "", existing)

		if len(targets) != 4 {
			t.Fatalf("Expected 4 targets (3 from md + 1 extra from existing), got %d", len(targets))
		}

		if targets[0].Title != "PAARA In The Park October 2026" || targets[0].FolderID != "1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82" {
			t.Errorf("Target 0 mismatch: %+v", targets[0])
		}
		if targets[2].Title != "Field Day 2019 Timelapse" || !targets[2].IsSingleFile {
			t.Errorf("Target 2 single file mismatch: %+v", targets[2])
		}
		if targets[3].Title != "Previously Indexed" || targets[3].FolderID != "previously-indexed-folder" {
			t.Errorf("Target 3 existing mismatch: %+v", targets[3])
		}
	})

	t.Run("filter by single folder with auto title detection", func(t *testing.T) {
		r := strings.NewReader(mdContent)
		targets := discoverTargetsFromReader(r, "1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82", "", nil)

		if len(targets) != 1 {
			t.Fatalf("Expected 1 target, got %d", len(targets))
		}
		if targets[0].Title != "PAARA In The Park October 2026" {
			t.Errorf("Expected title 'PAARA In The Park October 2026', got %q", targets[0].Title)
		}
	})

	t.Run("filter by single folder with existing map fallback", func(t *testing.T) {
		r := strings.NewReader("")
		targets := discoverTargetsFromReader(r, "previously-indexed-folder", "", existing)

		if len(targets) != 1 {
			t.Fatalf("Expected 1 target, got %d", len(targets))
		}
		if targets[0].Title != "Previously Indexed" {
			t.Errorf("Expected title from existing map, got %q", targets[0].Title)
		}
	})

	t.Run("filter by single folder with custom title override", func(t *testing.T) {
		r := strings.NewReader(mdContent)
		targets := discoverTargetsFromReader(r, "1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82", "Custom Field Day Title", nil)

		if len(targets) != 1 {
			t.Fatalf("Expected 1 target, got %d", len(targets))
		}
		if targets[0].Title != "Custom Field Day Title" {
			t.Errorf("Expected custom title override, got %q", targets[0].Title)
		}
	})

	t.Run("filter by single folder with default fallback title", func(t *testing.T) {
		r := strings.NewReader("")
		targets := discoverTargetsFromReader(r, "unknown-folder", "", nil)

		if len(targets) != 1 {
			t.Fatalf("Expected 1 target, got %d", len(targets))
		}
		if targets[0].Title != "Photo Album" {
			t.Errorf("Expected default 'Photo Album', got %q", targets[0].Title)
		}
	})
}

func TestDiscoverTargetsFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_photos.md")
	content := "## Disk Album\n{% include gallery.html folder_id=\"disk-folder-1\" %}\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	targets := discoverTargets(filePath, "", "", nil)
	if len(targets) != 1 || targets[0].FolderID != "disk-folder-1" {
		t.Errorf("discoverTargets from disk failed: %+v", targets)
	}

	// Test non-existent file
	targetsMissing := discoverTargets(filepath.Join(tempDir, "non_existent.md"), "single-id", "", nil)
	if len(targetsMissing) != 1 || targetsMissing[0].FolderID != "single-id" {
		t.Errorf("discoverTargets missing file fallback failed: %+v", targetsMissing)
	}
}

// ============================================================================
// Storage Round-Trip Tests
// ============================================================================

func TestLoadAndSaveGalleries(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "galleries.json")

	original := map[string]Gallery{
		"test-folder-1": {
			Title:      "Test Album 1",
			FolderID:   "test-folder-1",
			DriveURL:   "https://drive.google.com/drive/folders/test-folder-1",
			LastSynced: "2026-10-05T00:00:00Z",
			TotalCount: 1,
			ImageCount: 1,
			Items: []GalleryItem{
				{
					ID:        "item-1",
					Name:      "photo.jpg",
					Type:      "image",
					Timestamp: "2026-10-03T08:00:00",
					DateTime:  "Oct 3, 2026 at 8:00 AM",
				},
			},
		},
	}

	if err := saveGalleries(filePath, original); err != nil {
		t.Fatalf("saveGalleries failed: %v", err)
	}

	loaded := loadExistingGalleries(filePath)
	if len(loaded) != 1 || loaded["test-folder-1"].Title != "Test Album 1" {
		t.Errorf("Loaded gallery mismatch: %+v", loaded)
	}

	// Test load non-existent file returns empty map
	missing := loadExistingGalleries(filepath.Join(tempDir, "missing.json"))
	if missing == nil || len(missing) != 0 {
		t.Errorf("Expected empty map for missing file, got %+v", missing)
	}

	// Test save error with invalid path
	badPath := "/dev/null/impossible/dir/galleries.json"
	if err := saveGalleries(badPath, original); err == nil {
		t.Errorf("Expected error when saving to invalid directory path")
	}
}

// ============================================================================
// Mock HTTP Client
// ============================================================================

type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	urlStr := req.URL.String()
	if m.errors != nil {
		if err, ok := m.errors[urlStr]; ok && err != nil {
			return nil, err
		}
	}
	if m.responses != nil {
		if resp, ok := m.responses[urlStr]; ok {
			return resp, nil
		}
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}, nil
}

// ============================================================================
// Network Fetching Unit Tests
// ============================================================================

func TestFetchExifDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			responses: map[string]*http.Response{
				"https://lh3.googleusercontent.com/d/id-exif-ok=w1600": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("dummy...2026:10:03 08:30:00..."))),
				},
			},
		}

		iso, human := fetchExifDate(mockClient, "id-exif-ok")
		if iso != "2026-10-03T08:30:00" || human != "Oct 3, 2026 at 8:30 AM" {
			t.Errorf("fetchExifDate = (%q, %q); want (2026-10-03T08:30:00, Oct 3, 2026 at 8:30 AM)", iso, human)
		}
	})

	t.Run("client error returns empty", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			errors: map[string]error{
				"https://lh3.googleusercontent.com/d/id-exif-err=w1600": errors.New("network timeout"),
			},
		}
		iso, human := fetchExifDate(mockClient, "id-exif-err")
		if iso != "" || human != "" {
			t.Errorf("Expected empty strings on network error, got (%q, %q)", iso, human)
		}
	})
}

func TestFetchWithEmbeddedFolderView(t *testing.T) {
	embeddedURL := "https://drive.google.com/embeddedfolderview?id=folder-emb-1"

	t.Run("success", func(t *testing.T) {
		html := `<div class="flip-entry" id="entry-img1"><div class="flip-entry-title">photo.jpg</div></div>`
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embeddedURL: {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(html)),
				},
			},
		}
		items, err := fetchWithEmbeddedFolderView(client, "folder-emb-1", false)
		if err != nil || len(items) != 1 {
			t.Fatalf("fetchWithEmbeddedFolderView failed: items=%+v, err=%v", items, err)
		}
	})

	t.Run("HTTP 500 error", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embeddedURL: {
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("server error")),
				},
			},
		}
		_, err := fetchWithEmbeddedFolderView(client, "folder-emb-1", false)
		if err == nil {
			t.Fatalf("Expected error on HTTP 500")
		}
	})

	t.Run("empty entries returns error", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embeddedURL: {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("<div>No entries</div>")),
				},
			},
		}
		_, err := fetchWithEmbeddedFolderView(client, "folder-emb-1", false)
		if err == nil {
			t.Fatalf("Expected error when no entries found")
		}
	})

	t.Run("network error", func(t *testing.T) {
		client := &mockHTTPClient{
			errors: map[string]error{
				embeddedURL: errors.New("connection failed"),
			},
		}
		_, err := fetchWithEmbeddedFolderView(client, "folder-emb-1", false)
		if err == nil {
			t.Fatalf("Expected network error")
		}
	})
}

func TestFetchWithPublicScraper(t *testing.T) {
	publicURL := "https://drive.google.com/drive/folders/folder-pub-1"

	t.Run("success", func(t *testing.T) {
		html := `<div aria-label="photo.jpg Image" ssk='1:2:1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M-0-0'></div>`
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				publicURL: {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(html)),
				},
			},
		}
		items, err := fetchWithPublicScraper(client, "folder-pub-1", false)
		if err != nil || len(items) != 1 {
			t.Fatalf("fetchWithPublicScraper failed: items=%+v, err=%v", items, err)
		}
	})

	t.Run("HTTP 500 error", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				publicURL: {
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("server error")),
				},
			},
		}
		_, err := fetchWithPublicScraper(client, "folder-pub-1", false)
		if err == nil {
			t.Fatalf("Expected error on HTTP 500")
		}
	})

	t.Run("empty items returns error", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				publicURL: {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("<div>No aria-labels</div>")),
				},
			},
		}
		_, err := fetchWithPublicScraper(client, "folder-pub-1", false)
		if err == nil {
			t.Fatalf("Expected error when no items found")
		}
	})

	t.Run("network error", func(t *testing.T) {
		client := &mockHTTPClient{
			errors: map[string]error{
				publicURL: errors.New("connection failed"),
			},
		}
		_, err := fetchWithPublicScraper(client, "folder-pub-1", false)
		if err == nil {
			t.Fatalf("Expected network error")
		}
	})
}

func TestFetchWithDriveAPI(t *testing.T) {
	folderID := "folder-api-123"
	apiKey := "test-api-key"
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)

	page1URL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files?q=%s&fields=nextPageToken,files(id,name,mimeType,thumbnailLink,imageMediaMetadata,videoMediaMetadata)&orderBy=name&pageSize=1000&key=%s",
		url.QueryEscape(query), apiKey)

	page2URL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files?q=%s&fields=nextPageToken,files(id,name,mimeType,thumbnailLink,imageMediaMetadata,videoMediaMetadata)&orderBy=name&pageSize=1000&key=%s",
		url.QueryEscape(query), apiKey) + "&pageToken=page-token-2"

	t.Run("success with pagination, image metadata, and video", func(t *testing.T) {
		page1JSON := `{
			"nextPageToken": "page-token-2",
			"files": [
				{
					"id": "file-img-1",
					"name": "photo.jpg",
					"mimeType": "image/jpeg",
					"imageMediaMetadata": {
						"width": 1920,
						"height": 1080,
						"time": "2026:10:03 14:15:22"
					}
				}
			]
		}`

		page2JSON := `{
			"files": [
				{
					"id": "file-vid-1",
					"name": "clip.mp4",
					"mimeType": "video/mp4",
					"videoMediaMetadata": {
						"width": 1280,
						"height": 720
					}
				},
				{
					"id": "file-doc-1",
					"name": "notes.txt",
					"mimeType": "text/plain"
				}
			]
		}`

		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				page1URL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(page1JSON))},
				page2URL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(page2JSON))},
			},
		}

		items, err := fetchWithDriveAPI(client, folderID, apiKey, false)
		if err != nil {
			t.Fatalf("fetchWithDriveAPI failed: %v", err)
		}

		if len(items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(items))
		}

		if items[0].ID != "file-img-1" || items[0].Timestamp != "2026-10-03T14:15:22" || items[0].Width != 1920 {
			t.Errorf("Image item mismatch: %+v", items[0])
		}
		if items[1].ID != "file-vid-1" || items[1].Type != "video" || items[1].Width != 1280 {
			t.Errorf("Video item mismatch: %+v", items[1])
		}
	})

	t.Run("API error 403 response", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				page1URL: {StatusCode: http.StatusForbidden, Body: io.NopCloser(strings.NewReader("Rate limit exceeded"))},
			},
		}
		_, err := fetchWithDriveAPI(client, folderID, apiKey, false)
		if err == nil {
			t.Fatalf("Expected error on API 403 response")
		}
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				page1URL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("not json"))},
			},
		}
		_, err := fetchWithDriveAPI(client, folderID, apiKey, false)
		if err == nil {
			t.Fatalf("Expected error on invalid JSON")
		}
	})

	t.Run("client network error", func(t *testing.T) {
		client := &mockHTTPClient{
			errors: map[string]error{
				page1URL: errors.New("network unreachable"),
			},
		}
		_, err := fetchWithDriveAPI(client, folderID, apiKey, false)
		if err == nil {
			t.Fatalf("Expected network error")
		}
	})
}

func TestFetchSingleFile(t *testing.T) {
	singleURL := "https://lh3.googleusercontent.com/d/file-single=w400"

	t.Run("success video file", func(t *testing.T) {
		header := make(http.Header)
		header.Set("Content-Disposition", `attachment; filename="timelapse.mp4"`)
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				singleURL: {
					StatusCode: http.StatusOK,
					Header:     header,
					Body:       io.NopCloser(bytes.NewReader(nil)),
				},
			},
		}

		item, err := fetchSingleFile(client, "file-single", "Field Day Timelapse")
		if err != nil {
			t.Fatalf("fetchSingleFile failed: %v", err)
		}
		if item.Type != "video" || item.Name != "timelapse.mp4" {
			t.Errorf("Single video item mismatch: %+v", item)
		}
	})

	t.Run("success image file", func(t *testing.T) {
		header := make(http.Header)
		header.Set("Content-Disposition", `attachment; filename="group_photo.jpg"`)
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				singleURL: {
					StatusCode: http.StatusOK,
					Header:     header,
					Body:       io.NopCloser(bytes.NewReader(nil)),
				},
			},
		}

		item, err := fetchSingleFile(client, "file-single", "Group Photo")
		if err != nil {
			t.Fatalf("fetchSingleFile failed: %v", err)
		}
		if item.Type != "image" || item.Name != "group_photo.jpg" {
			t.Errorf("Single image item mismatch: %+v", item)
		}
	})

	t.Run("non-200 status code", func(t *testing.T) {
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				singleURL: {
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader(nil)),
				},
			},
		}
		_, err := fetchSingleFile(client, "file-single", "Group Photo")
		if err == nil {
			t.Fatalf("Expected error on HTTP 404")
		}
	})

	t.Run("client network error", func(t *testing.T) {
		client := &mockHTTPClient{
			errors: map[string]error{
				singleURL: errors.New("timeout"),
			},
		}
		_, err := fetchSingleFile(client, "file-single", "Group Photo")
		if err == nil {
			t.Fatalf("Expected network error")
		}
	})
}

// ============================================================================
// Enrichment & Pipeline Tests
// ============================================================================

func TestEnrichAndSortItems(t *testing.T) {
	t.Run("enrichment from filename, exif, and already set timestamp", func(t *testing.T) {
		mockClient := &mockHTTPClient{
			responses: map[string]*http.Response{
				"https://lh3.googleusercontent.com/d/id-exif-1=w1600": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("dummy...2026:10:03 08:30:00..."))),
				},
			},
		}

		items := []GalleryItem{
			{ID: "id-filename-1", Name: "20261003_140000.jpg", Type: "image"},
			{ID: "id-exif-1", Name: "IMG_0001.HEIC", Type: "image"},
			{ID: "id-cached-1", Name: "cached.jpg", Type: "image", Timestamp: "2026-10-03T10:00:00", DateTime: "Oct 3, 2026 at 10:00 AM"},
		}

		enriched := enrichAndSortItems(mockClient, items, nil, true)
		if len(enriched) != 3 {
			t.Fatalf("Expected 3 items, got %d", len(enriched))
		}

		expectedIDs := []string{"id-exif-1", "id-cached-1", "id-filename-1"}
		for i, it := range enriched {
			if it.ID != expectedIDs[i] {
				t.Errorf("Index %d: got ID %q, want %q", i, it.ID, expectedIDs[i])
			}
		}
	})

	t.Run("cached timestamp reuse from previous gallery sync", func(t *testing.T) {
		items := []GalleryItem{
			{ID: "id-from-cache", Name: "photo_without_date.jpg", Type: "image"},
		}
		existing := []GalleryItem{
			{ID: "id-from-cache", Name: "photo_without_date.jpg", Type: "image", Timestamp: "2026-05-01T10:00:00", DateTime: "May 1, 2026 at 10:00 AM"},
		}
		enriched := enrichAndSortItems(nil, items, existing, false)
		if len(enriched) != 1 || enriched[0].Timestamp != "2026-05-01T10:00:00" {
			t.Errorf("Expected timestamp reused from cache, got %+v", enriched)
		}
	})
}

func TestSyncTargetScraperFallbacks(t *testing.T) {
	target := SyncTarget{
		Title:    "Fallback Album",
		FolderID: "folder-fallback-1",
		DriveURL: "https://drive.google.com/drive/folders/folder-fallback-1",
	}

	t.Run("embedded scraper succeeds", func(t *testing.T) {
		embURL := "https://drive.google.com/embeddedfolderview?id=folder-fallback-1"
		html := `<div class="flip-entry" id="entry-f1"><div class="flip-entry-title">20261003_120000.jpg</div></div>`
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embURL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(html))},
			},
		}

		gal, err := syncTarget(client, target, Gallery{}, SyncOptions{})
		if err != nil || gal.TotalCount != 1 {
			t.Fatalf("syncTarget via embedded failed: gal=%+v, err=%v", gal, err)
		}
	})

	t.Run("embedded fails, public scraper succeeds", func(t *testing.T) {
		embURL := "https://drive.google.com/embeddedfolderview?id=folder-fallback-1"
		pubURL := "https://drive.google.com/drive/folders/folder-fallback-1"
		html := `<div aria-label="20261003_120000.jpg Image" ssk='1:2:1EELEkRrzcbpj4aV08QfiClkuXFQh3E0M-0-0'></div>`
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embURL: {StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("error"))},
				pubURL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(html))},
			},
		}

		gal, err := syncTarget(client, target, Gallery{}, SyncOptions{Verbose: true})
		if err != nil || gal.TotalCount != 1 {
			t.Fatalf("syncTarget via public scraper failed: gal=%+v, err=%v", gal, err)
		}
	})

	t.Run("scrapers fail, single file succeeds", func(t *testing.T) {
		singleURL := "https://lh3.googleusercontent.com/d/folder-fallback-1=w400"
		header := make(http.Header)
		header.Set("Content-Disposition", `attachment; filename="timelapse.mp4"`)
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				singleURL: {StatusCode: http.StatusOK, Header: header, Body: io.NopCloser(bytes.NewReader(nil))},
			},
		}

		gal, err := syncTarget(client, target, Gallery{}, SyncOptions{})
		if err != nil || gal.TotalCount != 1 || gal.VideoCount != 1 {
			t.Fatalf("syncTarget single file fallback failed: gal=%+v, err=%v", gal, err)
		}
	})

	t.Run("Drive API succeeds when key provided", func(t *testing.T) {
		query := fmt.Sprintf("'%s' in parents and trashed = false", target.FolderID)
		apiURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files?q=%s&fields=nextPageToken,files(id,name,mimeType,thumbnailLink,imageMediaMetadata,videoMediaMetadata)&orderBy=name&pageSize=1000&key=my-key",
			url.QueryEscape(query))
		jsonResp := `{"files":[{"id":"f-api-1","name":"20261003_120000.jpg","mimeType":"image/jpeg"}]}`

		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				apiURL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(jsonResp))},
			},
		}

		gal, err := syncTarget(client, target, Gallery{}, SyncOptions{APIKey: "my-key", Verbose: true})
		if err != nil || gal.TotalCount != 1 {
			t.Fatalf("syncTarget via Drive API failed: gal=%+v, err=%v", gal, err)
		}
	})

	t.Run("all fetchers fail, retains existing items", func(t *testing.T) {
		client := &mockHTTPClient{responses: map[string]*http.Response{}}
		existing := Gallery{
			Title:      "Existing",
			FolderID:   target.FolderID,
			TotalCount: 1,
			Items:      []GalleryItem{{ID: "old-1", Name: "photo.jpg", Type: "image"}},
		}

		gal, err := syncTarget(client, target, existing, SyncOptions{})
		if err != nil || gal.TotalCount != 1 {
			t.Fatalf("Expected retention of existing items: gal=%+v, err=%v", gal, err)
		}
	})

	t.Run("all fetchers fail, no existing returns error", func(t *testing.T) {
		client := &mockHTTPClient{responses: map[string]*http.Response{}}
		_, err := syncTarget(client, target, Gallery{}, SyncOptions{})
		if err == nil {
			t.Fatalf("Expected error when all fetchers fail and no existing gallery")
		}
	})

	t.Run("single file target", func(t *testing.T) {
		singleTarget := SyncTarget{
			Title:        "Single File Clip",
			FolderID:     "1single_file_id_test",
			DriveURL:     "https://drive.google.com/file/d/1single_file_id_test/view",
			IsSingleFile: true,
		}
		gal, err := syncTarget(nil, singleTarget, Gallery{}, SyncOptions{})
		if err != nil {
			t.Fatalf("syncTarget single file failed: %v", err)
		}
		if gal.TotalCount != 1 || gal.VideoCount != 1 || gal.ImageCount != 0 {
			t.Errorf("Unexpected single file gallery: %+v", gal)
		}
	})
}

// ============================================================================
// Run Orchestrator Tests
// ============================================================================

func TestRunWithClient(t *testing.T) {
	// Disable sleep during tests
	oldSleep := syncSleepDuration
	syncSleepDuration = 0
	defer func() { syncSleepDuration = oldSleep }()

	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "galleries.json")
	mdPath := filepath.Join(tempDir, "photos.md")

	t.Run("enrich-only mode with folder filter", func(t *testing.T) {
		initial := map[string]Gallery{
			"f1": {
				Title:    "Album 1",
				FolderID: "f1",
				Items: []GalleryItem{
					{ID: "1", Name: "20261003_120000.jpg", Type: "image"},
					{ID: "v1", Name: "clip.mp4", Type: "video", Timestamp: "2026-10-03T12:05:00"},
				},
			},
			"f2": {
				Title:    "Album 2",
				FolderID: "f2",
				Items:    []GalleryItem{{ID: "2", Name: "20261003_130000.jpg", Type: "image"}},
			},
		}
		if err := saveGalleries(jsonPath, initial); err != nil {
			t.Fatalf("saveGalleries failed: %v", err)
		}

		opts := SyncOptions{
			OutputFile: jsonPath,
			FolderID:   "f1",
			EnrichOnly: true,
		}
		if err := runWithClient(nil, opts); err != nil {
			t.Fatalf("runWithClient failed: %v", err)
		}

		loaded := loadExistingGalleries(jsonPath)
		if loaded["f1"].Items[0].DateTime != "Oct 3, 2026 at 12:00 PM" {
			t.Errorf("Expected f1 to be enriched, got %q", loaded["f1"].Items[0].DateTime)
		}
		if loaded["f1"].VideoCount != 1 || loaded["f1"].ImageCount != 1 {
			t.Errorf("Expected 1 image and 1 video, got %d images, %d videos", loaded["f1"].ImageCount, loaded["f1"].VideoCount)
		}
	})

	t.Run("enrich-only mode save error", func(t *testing.T) {
		opts := SyncOptions{
			OutputFile: "/dev/null/forbidden/path.json",
			EnrichOnly: true,
		}
		if err := runWithClient(nil, opts); err == nil {
			t.Fatalf("Expected error when saving fails in enrich-only mode")
		}
	})

	t.Run("full sync with discovered targets", func(t *testing.T) {
		syncSleepDuration = 1 * time.Millisecond
		defer func() { syncSleepDuration = 0 }()

		mdContent := "## Field Day\n{% include gallery.html folder_id=\"f-sync-1\" %}\n"
		if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}

		embURL := "https://drive.google.com/embeddedfolderview?id=f-sync-1"
		html := `<div class="flip-entry" id="entry-1"><div class="flip-entry-title">20261003_090000.jpg</div></div>`
		client := &mockHTTPClient{
			responses: map[string]*http.Response{
				embURL: {StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(html))},
			},
		}

		opts := SyncOptions{
			InputFile:  mdPath,
			OutputFile: jsonPath,
		}
		if err := runWithClient(client, opts); err != nil {
			t.Fatalf("runWithClient full sync failed: %v", err)
		}

		loaded := loadExistingGalleries(jsonPath)
		if len(loaded) == 0 || loaded["f-sync-1"].TotalCount != 1 {
			t.Errorf("Expected loaded gallery with 1 item: %+v", loaded)
		}
	})

	t.Run("full sync with no targets", func(t *testing.T) {
		emptyMdPath := filepath.Join(tempDir, "empty.md")
		if err := os.WriteFile(emptyMdPath, []byte("# Empty\n"), 0644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		opts := SyncOptions{
			InputFile:  emptyMdPath,
			OutputFile: filepath.Join(tempDir, "non_existent.json"),
		}
		if err := runWithClient(nil, opts); err != nil {
			t.Fatalf("Expected nil when no targets found, got %v", err)
		}
	})

	t.Run("full sync with target error and save failure", func(t *testing.T) {
		mdContent := "## Broken Folder\n{% include gallery.html folder_id=\"f-broken\" %}\n"
		if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}

		// Empty client (all requests 404)
		client := &mockHTTPClient{responses: map[string]*http.Response{}}
		opts := SyncOptions{
			InputFile:  mdPath,
			OutputFile: "/dev/null/forbidden/path.json",
		}
		if err := runWithClient(client, opts); err == nil {
			t.Fatalf("Expected save error on invalid output path")
		}
	})
}

func TestParseFlags(t *testing.T) {
	t.Run("default flags", func(t *testing.T) {
		opts, err := parseFlags([]string{})
		if err != nil {
			t.Fatalf("parseFlags default failed: %v", err)
		}
		if opts.InputFile != "photos.md" || opts.OutputFile != filepath.Join("_data", "galleries.json") {
			t.Errorf("Unexpected default options: %+v", opts)
		}
	})

	t.Run("all custom flags", func(t *testing.T) {
		args := []string{
			"-key=my-test-key",
			"-input=custom_input.md",
			"-output=custom_output.json",
			"-folder=custom_folder",
			"-title=Custom Title",
			"-enrich-only=true",
			"-verbose=true",
		}
		opts, err := parseFlags(args)
		if err != nil {
			t.Fatalf("parseFlags custom failed: %v", err)
		}
		if opts.APIKey != "my-test-key" || opts.InputFile != "custom_input.md" ||
			opts.OutputFile != "custom_output.json" || opts.FolderID != "custom_folder" ||
			opts.CustomTitle != "Custom Title" || !opts.EnrichOnly || !opts.Verbose {
			t.Errorf("Unexpected parsed options: %+v", opts)
		}
	})

	t.Run("invalid flag", func(t *testing.T) {
		_, err := parseFlags([]string{"-invalid-flag-xyz"})
		if err == nil {
			t.Errorf("Expected error for invalid flag, got nil")
		}
	})
}

func TestRun(t *testing.T) {
	tempDir := t.TempDir()
	emptyMdPath := filepath.Join(tempDir, "empty.md")
	if err := os.WriteFile(emptyMdPath, []byte("# Empty\n"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	opts := SyncOptions{
		InputFile:  emptyMdPath,
		OutputFile: filepath.Join(tempDir, "output.json"),
	}
	if err := Run(opts); err != nil {
		t.Errorf("Run() failed: %v", err)
	}
}
