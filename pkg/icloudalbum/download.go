// ABOUTME: Downloads photos from iCloud with automatic format detection and filename sanitization
// ABOUTME: Creates safe, descriptive filenames using GUIDs, captions, and indices
package icloudalbum

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// DownloadPhoto downloads the best derivative and writes it to outputDir.
// If customFilename (without extension) is given, it is used; otherwise we
// compose a name from GUID, caption (sanitized), and optional index.
func DownloadPhoto(photo *Image, index *int, outputDir string, customFilename *string) (string, error) {
	client := &http.Client{}

	key, _, url, ok := SelectBestDerivative(photo.Derivatives)
	if !ok || url == "" {
		return "", fmt.Errorf("no suitable derivative found (key=%q)", key)
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download failed (status %d)", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ext := GetExtensionForContent(content, "")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}

	base := ""
	switch {
	case customFilename != nil && *customFilename != "":
		base = fmt.Sprintf("%s_%s", photo.PhotoGUID, sanitize(*customFilename))
	case photo.Caption != nil && *photo.Caption != "":
		if index != nil {
			base = fmt.Sprintf("%d_%s_%s", *index+1, photo.PhotoGUID, sanitize(*photo.Caption))
		} else {
			base = fmt.Sprintf("%s_%s", photo.PhotoGUID, sanitize(*photo.Caption))
		}
	case index != nil:
		base = fmt.Sprintf("%d_%s", *index+1, photo.PhotoGUID)
	default:
		base = photo.PhotoGUID
	}

	fp := filepath.Join(outputDir, base+ext)
	if err := os.WriteFile(fp, content, 0o644); err != nil {
		return "", err
	}
	return fp, nil
}

func sanitize(s string) string {
	// Conservative, cross-platform safe: replace forbidden/awkward chars with '_'
	repl := func(r rune) rune {
		switch r {
		case '<', '>', ':', '"', '/', '\\', '|', '?', '*', '!', '@', '#', '$', '%', '^', '&', '\'', ';', '=', '+', ',', '`', '~':
			return '_'
		}
		// Strip control runes
		if r < 32 || r == 127 || !utf8.ValidRune(r) {
			return '_'
		}
		return r
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		out = append(out, repl(r))
	}
	if len(out) > 200 {
		out = append(out[:195], []rune("_truncated")...)
	}
	trimmed := string(out)
	trimmed = trimDots(trimmed)
	return trimmed
}

func trimDots(s string) string {
	// remove leading/trailing dots and whitespace
	for len(s) > 0 && (s[0] == '.' || s[0] == ' ') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == '.' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}
