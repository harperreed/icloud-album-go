// ABOUTME: Main orchestrator for fetching iCloud shared album data
// ABOUTME: Handles base URL calculation, redirects, API calls, and photo enrichment
package icloudalbum

import (
	"net/http"
	"time"
)

var defaultClient = &http.Client{
	Timeout: 30 * time.Second,
}

// GetICloudPhotos orchestrates:
// 1) base URL from token
// 2) 330 redirect handling
// 3) webstream metadata+photos
// 4) webasseturls URLs
// 5) enrichment of derivatives with URLs
func GetICloudPhotos(token string) (*ICloudResponse, error) {
	return GetICloudPhotosWithClient(token, defaultClient)
}

// GetICloudPhotosWithClient allows using a custom HTTP client for advanced use cases.
func GetICloudPhotosWithClient(token string, client *http.Client) (*ICloudResponse, error) {

	base, err := GetBaseURL(token)
	if err != nil {
		return nil, err
	}
	redirected, err := GetRedirectedBaseURL(client, base, token)
	if err != nil {
		return nil, err
	}

	photos, md, err := getAPIResponse(client, redirected)
	if err != nil {
		return nil, err
	}

	guids := make([]string, 0, len(photos))
	for _, p := range photos {
		guids = append(guids, p.PhotoGUID)
	}
	allURLs, err := GetAssetURLs(client, redirected, guids, nil)
	if err != nil {
		// Match Rust behavior: partial degradation is fine (e.g., 400 → empty map)
		// So we don't fail hard here; we just enrich with whatever we got.
	}

	EnrichPhotosWithURLs(photos, allURLs)

	return &ICloudResponse{
		Metadata: md,
		Photos:   photos,
	}, nil
}
