// ABOUTME: Utility functions for MIME type detection and derivative selection
// ABOUTME: Includes magic number detection for common image/video formats
package icloudalbum

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// ExtensionFromMIME maps a MIME type to a file extension.
func ExtensionFromMIME(mt string) string {
	switch mt {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/heic":
		return ".heic"
	case "image/heif":
		return ".heif"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "image/gif":
		return ".gif"
	default:
		log.Printf("warn: unknown MIME type %q; defaulting to .jpg", mt)
		return ".jpg"
	}
}

// DetectMIMEType inspects bytes (and optionally filename) to guess the MIME type.
func DetectMIMEType(b []byte, filename string) string {
	// Manual magic numbers (parity with Rust)
	if len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF {
		return "image/jpeg"
	}
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return "image/png"
	}
	if len(b) > 11 && b[4] == 0x66 && b[5] == 0x74 && b[6] == 0x79 && b[7] == 0x70 {
		// 'ftyp'
		if b[8] == 0x71 && b[9] == 0x74 { // 'qt'
			return "video/quicktime"
		}
		return "video/mp4"
	}
	if len(b) >= 6 &&
		b[0] == 0x47 && b[1] == 0x49 && b[2] == 0x46 && b[3] == 0x38 &&
		(b[4] == 0x37 || b[4] == 0x39) && b[5] == 0x61 {
		return "image/gif"
	}
	if len(b) > 12 &&
		b[4] == 0x66 && b[5] == 0x74 && b[6] == 0x79 && b[7] == 0x70 &&
		b[8] == 0x68 && b[9] == 0x65 && b[10] == 0x69 &&
		(b[11] == 0x63 || b[11] == 0x66) {
		if b[11] == 0x63 {
			return "image/heic"
		}
		return "image/heif"
	}

	// Fallback to sniffing
	mt := http.DetectContentType(b)
	if mt == "application/octet-stream" && filename != "" {
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != "" {
			if t := mime.TypeByExtension(ext); t != "" {
				return t
			}
		}
	}
	if mt == "" {
		return "image/jpeg"
	}
	return mt
}

func GetExtensionForContent(b []byte, filename string) string {
	return ExtensionFromMIME(DetectMIMEType(b, filename))
}

// SelectBestDerivative mirrors the Rust logic.
// 1) Prefer originals ("original", "full", keys "3" or "4") with dimensions.
// 2) Otherwise highest resolution with dimensions.
// 3) Otherwise first derivative that has a URL.
func SelectBestDerivative(derivs map[string]Derivative) (key string, d Derivative, url string, ok bool) {
	if len(derivs) == 0 {
		return "", Derivative{}, "", false
	}

	var (
		bestKey   string
		best      *Derivative
		maxPixels uint64
		haveOrig  bool
	)

	for k, v := range derivs {
		if v.URL == nil {
			continue
		}
		isOriginal := strings.Contains(strings.ToLower(k), "original") ||
			strings.Contains(strings.ToLower(k), "full") ||
			k == "3" || k == "4"

		if isOriginal {
			haveOrig = true
			if v.Width != nil && v.Height != nil {
				pix := uint64(*v.Width) * uint64(*v.Height)
				if pix > maxPixels {
					maxPixels = pix
					bestKey, best = k, &v
				}
			} else if best == nil {
				// No dimensions but still prefer some original if nothing better yet
				bestKey, best = k, &v
			}
			continue
		}

		// Non-original candidates if we don't yet have a good original
		if !haveOrig && v.Width != nil && v.Height != nil {
			pix := uint64(*v.Width) * uint64(*v.Height)
			if pix > maxPixels {
				maxPixels = pix
				bestKey, best = k, &v
			}
		}
	}

	if best != nil {
		return bestKey, *best, *best.URL, true
	}

	// Fallback: first URL
	for k, v := range derivs {
		if v.URL != nil {
			return k, v, *v.URL, true
		}
	}
	return "", Derivative{}, "", false
}
