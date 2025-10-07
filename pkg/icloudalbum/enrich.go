// ABOUTME: Enriches photo derivatives with URLs from asset URL mappings
// ABOUTME: Uses lenient strategy matching by checksum first, then photo GUID
package icloudalbum

// EnrichPhotosWithURLs populates derivative URLs using a lenient strategy:
// 1) If a derivative checksum matches a key in allURLs, use that.
// 2) Otherwise, if the photo GUID matches a key in allURLs, apply that URL
//    to any derivative missing a URL.
// This mirrors the mixed real-world behavior (sometimes Apple keys by checksum,
// sometimes by photo GUID) and preserves your Rust tests' intent.
func EnrichPhotosWithURLs(photos []Image, allURLs map[string]string) {
	for pi := range photos {
		p := &photos[pi]

		// First pass: try checksum matches
		for k, d := range p.Derivatives {
			if d.URL == nil {
				if u, ok := allURLs[d.Checksum]; ok {
					url := u
					d.URL = &url
					p.Derivatives[k] = d
				}
			}
		}
		// Second pass: fallback to per-photo GUID
		if u, ok := allURLs[p.PhotoGUID]; ok {
			for k, d := range p.Derivatives {
				if d.URL == nil {
					url := u
					d.URL = &url
					p.Derivatives[k] = d
				}
			}
		}
	}
}
