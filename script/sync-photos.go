// script/sync-photos.go
// Synchronizes Google Drive photo folders referenced in photos.md and event pages
// into _data/galleries.json for embedding responsive photo galleries on the PAARA website.
//
// Usage:
//   go run ./script/sync-photos.go
//   go run ./script/sync-photos.go -folder="1ccgA24vYa5udg7r3pI-JyGQXdU1KvO82"
//   go run ./script/sync-photos.go -enrich-only
//
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// HTTPClient defines the interface needed for performing HTTP requests.
// This allows easy mocking in unit tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GalleryItem represents an individual photo or video.
type GalleryItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"` // "image" or "video"
	Thumb     string `json:"thumb"`
	Full      string `json:"full,omitempty"`      // For images
	Preview   string `json:"preview,omitempty"`   // For videos (Google Drive preview player)
	Timestamp string `json:"timestamp,omitempty"` // ISO 8601 sortable timestamp (e.g. 2026-06-13T05:56:14)
	DateTime  string `json:"date_time,omitempty"` // Human-readable date & time (e.g. Jun 13, 2026 at 5:56 AM)
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

// Gallery represents an event's media album.
type Gallery struct {
	Title      string        `json:"title"`
	FolderID   string        `json:"folder_id"`
	DriveURL   string        `json:"drive_url"`
	LastSynced string        `json:"last_synced"`
	TotalCount int           `json:"total_count"`
	ImageCount int           `json:"image_count"`
	VideoCount int           `json:"video_count"`
	Items      []GalleryItem `json:"items"`
}

// SyncTarget defines a Google Drive folder or file target to be synchronized.
type SyncTarget struct {
	Title        string
	FolderID     string
	DriveURL     string
	IsSingleFile bool
}

// SyncOptions holds the configuration options for a synchronization run.
type SyncOptions struct {
	APIKey      string
	InputFile   string
	OutputFile  string
	FolderID    string
	CustomTitle string
	EnrichOnly  bool
	Verbose     bool
}

// DriveAPIFileList maps Google Drive API v3 files list responses.
type DriveAPIFileList struct {
	NextPageToken string `json:"nextPageToken"`
	Files         []struct {
		ID                 string `json:"id"`
		Name               string `json:"name"`
		MimeType           string `json:"mimeType"`
		ThumbnailLink      string `json:"thumbnailLink"`
		ImageMediaMetadata *struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Time   string `json:"time"`
		} `json:"imageMediaMetadata"`
		VideoMediaMetadata *struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"videoMediaMetadata"`
	} `json:"files"`
}

var (
	folderRegex         = regexp.MustCompile(`https://drive\.google\.com/drive/folders/([a-zA-Z0-9_-]+)`)
	fileRegex           = regexp.MustCompile(`https://drive\.google\.com/file/d/([a-zA-Z0-9_-]+)`)
	galleryIncludeRegex = regexp.MustCompile(`folder_id=["']([a-zA-Z0-9_-]+)["']`)
	titleRegex          = regexp.MustCompile(`^##+\s+(.+)$`)
	sessionSuffixRegex  = regexp.MustCompile(`-\d+-\d+$`)

	// Scrapers for public folder pages
	embeddedEntryRegex = regexp.MustCompile(`class="flip-entry"[^>]*id="entry-([a-zA-Z0-9_-]+)"[\s\S]*?<div class="flip-entry-title">([^<]+)</div>`)
	ariaLabelRegex     = regexp.MustCompile(`aria-label="([^"]+)\s+(Image|Video)"[^>]*ssk='[^':]*:[^':]*:([a-zA-Z0-9_-]{25,})`)

	// Timestamp regexes
	exifDateRegex     = regexp.MustCompile(`\b(19\d\d|20\d\d):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2})\b`)
	filenameDateRegex = regexp.MustCompile(`(?:^|[^0-9])(20\d\d)(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})(?:[^0-9]|$)`)
)

// ============================================================================
// Classification & Pure Parsing Functions
// ============================================================================

// isImageFile returns true if the filename has a known image file extension.
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".bmp", ".tif", ".tiff":
		return true
	default:
		return false
	}
}

// isVideoFile returns true if the filename has a known video file extension.
func isVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".mov", ".m4v", ".webm", ".avi", ".mkv":
		return true
	default:
		return false
	}
}

// parseFilenameDate extracts ISO and human-readable timestamps from filenames
// formatted like YYYYMMDD_HHMMSS (e.g. 20261003_142311.jpg).
func parseFilenameDate(name string) (iso, human string) {
	m := filenameDateRegex.FindStringSubmatch(name)
	if len(m) < 7 {
		return "", ""
	}
	raw := fmt.Sprintf("%s-%s-%s %s:%s:%s", m[1], m[2], m[3], m[4], m[5], m[6])
	t, err := time.Parse("2006-01-02 15:04:05", raw)
	if err != nil {
		return "", ""
	}
	return t.Format("2006-01-02T15:04:05"), t.Format("Jan 2, 2006 at 3:04 PM")
}

// parseExifDate scans a byte buffer for standard EXIF date format (YYYY:MM:DD HH:MM:SS)
// and returns both ISO and human-readable representations.
func parseExifDate(buf []byte) (iso, human string) {
	if len(buf) == 0 {
		return "", ""
	}
	if m := exifDateRegex.FindSubmatch(buf); len(m) > 0 {
		raw := string(m[0])
		t, err := time.Parse("2006:01:02 15:04:05", raw)
		if err == nil {
			return t.Format("2006-01-02T15:04:05"), t.Format("Jan 2, 2006 at 3:04 PM")
		}
	}
	return "", ""
}

// parseEmbeddedFolderViewHTML extracts gallery items from the HTML of an embeddedfolderview page.
func parseEmbeddedFolderViewHTML(html string) []GalleryItem {
	matches := embeddedEntryRegex.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var items []GalleryItem

	for _, m := range matches {
		id := m[1]
		name := strings.TrimSpace(m[2])

		if seen[id] {
			continue
		}
		seen[id] = true

		item := GalleryItem{
			ID:    id,
			Name:  name,
			Thumb: fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w400", id),
		}

		if isVideoFile(name) {
			item.Type = "video"
			item.Preview = fmt.Sprintf("https://drive.google.com/file/d/%s/preview", id)
		} else {
			item.Type = "image"
			item.Full = fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w1600", id)
		}

		items = append(items, item)
	}

	sortGalleryItems(items)
	return items
}

// parsePublicScraperHTML extracts gallery items from the public Drive folder SPA HTML.
func parsePublicScraperHTML(html string) []GalleryItem {
	matches := ariaLabelRegex.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var items []GalleryItem

	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		mediaKind := strings.ToLower(m[2])
		rawID := m[3]

		cleanID := sessionSuffixRegex.ReplaceAllString(rawID, "")
		if seen[cleanID] {
			continue
		}
		seen[cleanID] = true

		item := GalleryItem{
			ID:    cleanID,
			Name:  name,
			Thumb: fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w400", cleanID),
		}

		if mediaKind == "video" || isVideoFile(name) {
			item.Type = "video"
			item.Preview = fmt.Sprintf("https://drive.google.com/file/d/%s/preview", cleanID)
		} else {
			item.Type = "image"
			item.Full = fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w1600", cleanID)
		}

		items = append(items, item)
	}

	sortGalleryItems(items)
	return items
}

// sortGalleryItems sorts gallery items chronologically ascending by timestamp,
// placing items with timestamps first, followed by alphabetical order by name.
func sortGalleryItems(items []GalleryItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Timestamp != "" && items[j].Timestamp != "" {
			if items[i].Timestamp != items[j].Timestamp {
				return items[i].Timestamp < items[j].Timestamp
			}
		} else if items[i].Timestamp != "" {
			return true
		} else if items[j].Timestamp != "" {
			return false
		}
		return items[i].Name < items[j].Name
	})
}

// discoverTargetsFromReader parses a markdown reader to find Google Drive gallery targets.
func discoverTargetsFromReader(r io.Reader, onlyFolder, customTitle string, existing map[string]Gallery) []SyncTarget {
	if onlyFolder != "" {
		title := customTitle
		if title == "" {
			if ex, ok := existing[onlyFolder]; ok && ex.Title != "" {
				title = ex.Title
			}
		}
		if title == "" && r != nil {
			currentTitle := ""
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if m := titleRegex.FindStringSubmatch(line); len(m) > 1 {
					currentTitle = strings.TrimSpace(m[1])
				}
				if strings.Contains(line, onlyFolder) && currentTitle != "" {
					title = currentTitle
					break
				}
			}
		}
		if title == "" {
			title = "Photo Album"
		}

		return []SyncTarget{{
			Title:    title,
			FolderID: onlyFolder,
			DriveURL: "https://drive.google.com/drive/folders/" + onlyFolder,
		}}
	}

	var targets []SyncTarget
	seen := make(map[string]bool)

	if r != nil {
		currentTitle := "Photo Album"
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if m := titleRegex.FindStringSubmatch(line); len(m) > 1 {
				currentTitle = strings.TrimSpace(m[1])
			}

			if m := galleryIncludeRegex.FindStringSubmatch(line); len(m) > 1 {
				fID := m[1]
				if !seen[fID] {
					seen[fID] = true
					targets = append(targets, SyncTarget{
						Title:        currentTitle,
						FolderID:     fID,
						DriveURL:     "https://drive.google.com/drive/folders/" + fID,
						IsSingleFile: false,
					})
				}
			} else if m := folderRegex.FindStringSubmatch(line); len(m) > 1 {
				fID := m[1]
				if !seen[fID] {
					seen[fID] = true
					targets = append(targets, SyncTarget{
						Title:        currentTitle,
						FolderID:     fID,
						DriveURL:     "https://drive.google.com/drive/folders/" + fID,
						IsSingleFile: false,
					})
				}
			} else if m := fileRegex.FindStringSubmatch(line); len(m) > 1 {
				fID := m[1]
				if !seen[fID] {
					seen[fID] = true
					targets = append(targets, SyncTarget{
						Title:        currentTitle,
						FolderID:     fID,
						DriveURL:     "https://drive.google.com/file/d/" + fID + "/view",
						IsSingleFile: true,
					})
				}
			}
		}
	}

	// Also include any previously indexed galleries from existing map
	for fID, g := range existing {
		if !seen[fID] {
			seen[fID] = true
			targets = append(targets, SyncTarget{
				Title:        g.Title,
				FolderID:     fID,
				DriveURL:     g.DriveURL,
				IsSingleFile: g.TotalCount == 1 && g.VideoCount == 1,
			})
		}
	}

	return targets
}

func discoverTargets(filename, onlyFolder, customTitle string, existing map[string]Gallery) []SyncTarget {
	var r io.Reader
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		r = file
		return discoverTargetsFromReader(r, onlyFolder, customTitle, existing)
	}
	return discoverTargetsFromReader(nil, onlyFolder, customTitle, existing)
}

// ============================================================================
// Network Fetching & EXIF Enrichment
// ============================================================================

// fetchExifDate reads up to the first 64KB of an image via HTTP Range request
// and parses the EXIF capture date.
func fetchExifDate(client HTTPClient, fileID string) (string, string) {
	urlStr := fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w1600", fileID)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", ""
	}
	req.Header.Set("Range", "bytes=0-65535")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	buf := make([]byte, 65536)
	n, _ := io.ReadFull(resp.Body, buf)
	return parseExifDate(buf[:n])
}

// fetchWithEmbeddedFolderView uses Google Drive's embeddedfolderview endpoint.
func fetchWithEmbeddedFolderView(client HTTPClient, folderID string, verbose bool) ([]GalleryItem, error) {
	urlStr := fmt.Sprintf("https://drive.google.com/embeddedfolderview?id=%s", folderID)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	items := parseEmbeddedFolderViewHTML(string(body))
	if len(items) == 0 {
		return nil, fmt.Errorf("no entries found in embeddedfolderview")
	}
	return items, nil
}

// fetchWithPublicScraper falls back to scraping the public Google Drive folder page.
func fetchWithPublicScraper(client HTTPClient, folderID string, verbose bool) ([]GalleryItem, error) {
	urlStr := fmt.Sprintf("https://drive.google.com/drive/folders/%s", folderID)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	items := parsePublicScraperHTML(string(body))
	if len(items) == 0 {
		return nil, fmt.Errorf("no media items found in public folder HTML")
	}
	return items, nil
}

// fetchWithDriveAPI fetches media items using Google Drive API v3.
func fetchWithDriveAPI(client HTTPClient, folderID, apiKey string, verbose bool) ([]GalleryItem, error) {
	var items []GalleryItem
	pageToken := ""

	for {
		query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
		reqURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files?q=%s&fields=nextPageToken,files(id,name,mimeType,thumbnailLink,imageMediaMetadata,videoMediaMetadata)&orderBy=name&pageSize=1000&key=%s",
			url.QueryEscape(query), apiKey)

		if pageToken != "" {
			reqURL += "&pageToken=" + url.QueryEscape(pageToken)
		}

		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
		}

		var fileList DriveAPIFileList
		if err := json.NewDecoder(resp.Body).Decode(&fileList); err != nil {
			return nil, err
		}

		for _, f := range fileList.Files {
			isImg := strings.HasPrefix(f.MimeType, "image/") || isImageFile(f.Name)
			isVid := strings.HasPrefix(f.MimeType, "video/") || isVideoFile(f.Name)

			if !isImg && !isVid {
				continue
			}

			item := GalleryItem{
				ID:    f.ID,
				Name:  f.Name,
				Thumb: fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w400", f.ID),
			}

			if isVid {
				item.Type = "video"
				item.Preview = fmt.Sprintf("https://drive.google.com/file/d/%s/preview", f.ID)
				if f.VideoMediaMetadata != nil {
					item.Width = f.VideoMediaMetadata.Width
					item.Height = f.VideoMediaMetadata.Height
				}
			} else {
				item.Type = "image"
				item.Full = fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w1600", f.ID)
				if f.ImageMediaMetadata != nil {
					item.Width = f.ImageMediaMetadata.Width
					item.Height = f.ImageMediaMetadata.Height
					if f.ImageMediaMetadata.Time != "" {
						if t, err := time.Parse("2006:01:02 15:04:05", f.ImageMediaMetadata.Time); err == nil {
							item.Timestamp = t.Format("2006-01-02T15:04:05")
							item.DateTime = t.Format("Jan 2, 2006 at 3:04 PM")
						}
					}
				}
			}

			items = append(items, item)
		}

		pageToken = fileList.NextPageToken
		if pageToken == "" {
			break
		}
	}

	sortGalleryItems(items)
	return items, nil
}

// fetchSingleFile verifies and registers a single Google Drive file (e.g. timelapse video).
func fetchSingleFile(client HTTPClient, fileID, title string) (GalleryItem, error) {
	urlStr := fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w400", fileID)
	req, err := http.NewRequest("HEAD", urlStr, nil)
	if err != nil {
		return GalleryItem{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return GalleryItem{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GalleryItem{}, fmt.Errorf("not a single file (HTTP %d)", resp.StatusCode)
	}

	filename := title
	disp := resp.Header.Get("Content-Disposition")
	if m := regexp.MustCompile(`filename="?([^";]+)"?`).FindStringSubmatch(disp); len(m) > 1 {
		filename = m[1]
	}

	item := GalleryItem{
		ID:    fileID,
		Name:  filename,
		Thumb: urlStr,
	}

	if isVideoFile(filename) || isVideoFile(title) || strings.Contains(strings.ToLower(title), "timelapse") {
		item.Type = "video"
		item.Preview = fmt.Sprintf("https://drive.google.com/file/d/%s/preview", fileID)
	} else {
		item.Type = "image"
		item.Full = fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w1600", fileID)
	}

	return item, nil
}

// enrichAndSortItems enriches items with timestamps (from cache, filename, or EXIF)
// using a concurrent worker pool, and sorts them chronologically.
func enrichAndSortItems(client HTTPClient, items []GalleryItem, existingItems []GalleryItem, verbose bool) []GalleryItem {
	cached := make(map[string]GalleryItem)
	for _, it := range existingItems {
		cached[it.ID] = it
	}

	var fetchIndices []int

	for i := range items {
		if items[i].Timestamp != "" {
			continue
		}

		// 1. Reuse existing timestamp from cache if available
		if ex, ok := cached[items[i].ID]; ok && ex.Timestamp != "" {
			items[i].Timestamp = ex.Timestamp
			items[i].DateTime = ex.DateTime
			continue
		}

		// 2. Check filename date
		if iso, human := parseFilenameDate(items[i].Name); iso != "" {
			items[i].Timestamp = iso
			items[i].DateTime = human
			continue
		}

		// 3. Queue for EXIF fetch if image
		if items[i].Type == "image" {
			fetchIndices = append(fetchIndices, i)
		}
	}

	if len(fetchIndices) > 0 {
		if verbose || len(fetchIndices) > 0 {
			fmt.Printf("    Extracting EXIF dates for %d items via Range requests...\n", len(fetchIndices))
		}
		numWorkers := 20
		if len(fetchIndices) < numWorkers {
			numWorkers = len(fetchIndices)
		}

		jobChan := make(chan int, len(fetchIndices))
		for _, idx := range fetchIndices {
			jobChan <- idx
		}
		close(jobChan)

		var wg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for idx := range jobChan {
					iso, human := fetchExifDate(client, items[idx].ID)
					if iso != "" {
						items[idx].Timestamp = iso
						items[idx].DateTime = human
					}
				}
			}()
		}
		wg.Wait()
	}

	sortGalleryItems(items)
	return items
}

// ============================================================================
// Storage & File I/O
// ============================================================================

// loadExistingGalleries loads previously synced gallery data from galleries.json.
func loadExistingGalleries(path string) map[string]Gallery {
	galleries := make(map[string]Gallery)
	data, err := os.ReadFile(path)
	if err != nil {
		return galleries
	}
	_ = json.Unmarshal(data, &galleries)
	return galleries
}

// saveGalleries writes gallery data to galleries.json formatted with indent.
func saveGalleries(path string, galleries map[string]Gallery) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(galleries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ============================================================================
// Sync Pipeline & Execution
// ============================================================================

// syncTarget orchestrates fetching and enriching a single target.
func syncTarget(client HTTPClient, target SyncTarget, existingGallery Gallery, opts SyncOptions) (Gallery, error) {
	if target.IsSingleFile {
		item := GalleryItem{
			ID:      target.FolderID,
			Name:    target.Title,
			Type:    "video",
			Thumb:   fmt.Sprintf("https://lh3.googleusercontent.com/d/%s=w400", target.FolderID),
			Preview: fmt.Sprintf("https://drive.google.com/file/d/%s/preview", target.FolderID),
		}
		return Gallery{
			Title:      target.Title,
			FolderID:   target.FolderID,
			DriveURL:   target.DriveURL,
			LastSynced: time.Now().UTC().Format(time.RFC3339),
			TotalCount: 1,
			ImageCount: 0,
			VideoCount: 1,
			Items:      []GalleryItem{item},
		}, nil
	}

	var items []GalleryItem
	var err error

	// 1. Try Google Drive API v3 if API key provided
	if opts.APIKey != "" {
		if opts.Verbose {
			fmt.Println("  Attempting Google Drive API v3...")
		}
		items, err = fetchWithDriveAPI(client, target.FolderID, opts.APIKey, opts.Verbose)
	}

	// 2. Use embedded folderview scraper
	if len(items) == 0 {
		items, err = fetchWithEmbeddedFolderView(client, target.FolderID, opts.Verbose)
	}

	// 3. Fallback to public folder SPA scraper
	if err != nil || len(items) == 0 {
		if opts.Verbose {
			fmt.Println("  Embedded view empty, falling back to public SPA scraper...")
		}
		items, err = fetchWithPublicScraper(client, target.FolderID, opts.Verbose)
	}

	// 4. Fallback: single file check
	if len(items) == 0 {
		if singleItem, singleErr := fetchSingleFile(client, target.FolderID, target.Title); singleErr == nil {
			items = []GalleryItem{singleItem}
			err = nil
		}
	}

	// 5. Retain existing items if network error occurs
	if err != nil || len(items) == 0 {
		if len(existingGallery.Items) > 0 {
			fmt.Printf("  Retaining %d existing item(s) from previous sync.\n", existingGallery.TotalCount)
			items = existingGallery.Items
			err = nil
		} else {
			return Gallery{}, fmt.Errorf("failed to fetch folder %s: %w", target.FolderID, err)
		}
	}

	// 6. Enrich with EXIF dates & sort chronologically
	items = enrichAndSortItems(client, items, existingGallery.Items, opts.Verbose)

	imgCount, vidCount := 0, 0
	for _, it := range items {
		if it.Type == "video" {
			vidCount++
		} else {
			imgCount++
		}
	}

	return Gallery{
		Title:      target.Title,
		FolderID:   target.FolderID,
		DriveURL:   target.DriveURL,
		LastSynced: time.Now().UTC().Format(time.RFC3339),
		TotalCount: len(items),
		ImageCount: imgCount,
		VideoCount: vidCount,
		Items:      items,
	}, nil
}

var syncSleepDuration = 250 * time.Millisecond

// Run executes the synchronization process with the provided options.
func Run(opts SyncOptions) error {
	client := &http.Client{Timeout: 30 * time.Second}
	return runWithClient(client, opts)
}

func runWithClient(client HTTPClient, opts SyncOptions) error {
	galleries := loadExistingGalleries(opts.OutputFile)

	// Mode 1: Enrich-only (no scraping)
	if opts.EnrichOnly {
		fmt.Printf("Enriching %d existing gallerie(s) in %s with EXIF timestamps...\n", len(galleries), opts.OutputFile)
		for folderID, gal := range galleries {
			if opts.FolderID != "" && folderID != opts.FolderID {
				continue
			}
			fmt.Printf("Enriching %s (%d items)...\n", gal.Title, len(gal.Items))
			gal.Items = enrichAndSortItems(client, gal.Items, gal.Items, opts.Verbose)
			gal.TotalCount = len(gal.Items)
			imgCount, vidCount := 0, 0
			for _, it := range gal.Items {
				if it.Type == "video" {
					vidCount++
				} else {
					imgCount++
				}
			}
			gal.ImageCount = imgCount
			gal.VideoCount = vidCount
			galleries[folderID] = gal
		}
		if err := saveGalleries(opts.OutputFile, galleries); err != nil {
			return fmt.Errorf("error saving %s: %w", opts.OutputFile, err)
		}
		fmt.Printf("Successfully updated and sorted all galleries in %s!\n", opts.OutputFile)
		return nil
	}

	// Mode 2: Full sync
	targets := discoverTargets(opts.InputFile, opts.FolderID, opts.CustomTitle, galleries)
	if len(targets) == 0 {
		fmt.Printf("No Google Drive targets found in %s\n", opts.InputFile)
		return nil
	}

	fmt.Printf("Found %d Google Drive target(s) to synchronize.\n", len(targets))

	updatedCount := 0
	for i, target := range targets {
		fmt.Printf("[%d/%d] Syncing: %s (ID: %s)...\n", i+1, len(targets), target.Title, target.FolderID)

		gal, err := syncTarget(client, target, galleries[target.FolderID], opts)
		if err != nil {
			fmt.Printf("  ⚠️ %v\n", err)
			continue
		}

		galleries[target.FolderID] = gal
		fmt.Printf("  ✓ Found %d media item(s) (%d photos, %d videos)\n", gal.TotalCount, gal.ImageCount, gal.VideoCount)
		updatedCount++

		if syncSleepDuration > 0 {
			time.Sleep(syncSleepDuration)
		}
	}

	if err := saveGalleries(opts.OutputFile, galleries); err != nil {
		return fmt.Errorf("error saving %s: %w", opts.OutputFile, err)
	}

	fmt.Printf("\nSuccessfully updated %s with %d galleries (%d updated in this run).\n", opts.OutputFile, len(galleries), updatedCount)
	return nil
}

// parseFlags parses command-line arguments into SyncOptions.
func parseFlags(args []string) (SyncOptions, error) {
	fs := flag.NewFlagSet("sync-photos", flag.ContinueOnError)
	var opts SyncOptions
	fs.StringVar(&opts.APIKey, "key", os.Getenv("GOOGLE_API_KEY"), "Google Drive API Key (optional)")
	fs.StringVar(&opts.InputFile, "input", "photos.md", "Markdown file to scan for Google Drive links")
	fs.StringVar(&opts.OutputFile, "output", filepath.Join("_data", "galleries.json"), "Output JSON file path")
	fs.StringVar(&opts.FolderID, "folder", "", "Sync only a specific folder ID")
	fs.StringVar(&opts.CustomTitle, "title", "", "Custom title when syncing a specific folder ID")
	fs.BoolVar(&opts.EnrichOnly, "enrich-only", false, "Only enrich existing galleries with EXIF timestamps and re-sort without re-scraping")
	fs.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose logging")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	return opts, nil
}

func main() {
	opts, err := parseFlags(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	if err := Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
