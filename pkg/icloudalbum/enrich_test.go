// ABOUTME: Test suite for photo URL enrichment logic
// ABOUTME: Validates checksum and GUID-based URL matching strategies
package icloudalbum

import (
	"testing"
)

func TestEnrichPhotosWithURLs(t *testing.T) {
	tests := []struct {
		name    string
		photos  []Image
		allURLs map[string]string
		verify  func(t *testing.T, photos []Image)
	}{
		{
			name: "enrich by checksum",
			photos: []Image{
				{
					PhotoGUID: "photo1",
					Derivatives: map[string]Derivative{
						"original": {Checksum: "checksum1"},
					},
				},
			},
			allURLs: map[string]string{
				"checksum1": "https://example.com/photo1.jpg",
			},
			verify: func(t *testing.T, photos []Image) {
				if photos[0].Derivatives["original"].URL == nil {
					t.Error("expected URL to be set")
					return
				}
				if *photos[0].Derivatives["original"].URL != "https://example.com/photo1.jpg" {
					t.Errorf("URL = %q, want https://example.com/photo1.jpg", *photos[0].Derivatives["original"].URL)
				}
			},
		},
		{
			name: "enrich by photo GUID fallback",
			photos: []Image{
				{
					PhotoGUID: "photo1",
					Derivatives: map[string]Derivative{
						"original": {Checksum: "unknown"},
					},
				},
			},
			allURLs: map[string]string{
				"photo1": "https://example.com/photo1.jpg",
			},
			verify: func(t *testing.T, photos []Image) {
				if photos[0].Derivatives["original"].URL == nil {
					t.Error("expected URL to be set")
					return
				}
				if *photos[0].Derivatives["original"].URL != "https://example.com/photo1.jpg" {
					t.Errorf("URL = %q, want https://example.com/photo1.jpg", *photos[0].Derivatives["original"].URL)
				}
			},
		},
		{
			name: "checksum takes priority over GUID",
			photos: []Image{
				{
					PhotoGUID: "photo1",
					Derivatives: map[string]Derivative{
						"original": {Checksum: "checksum1"},
					},
				},
			},
			allURLs: map[string]string{
				"checksum1": "https://example.com/checksum.jpg",
				"photo1":    "https://example.com/guid.jpg",
			},
			verify: func(t *testing.T, photos []Image) {
				if photos[0].Derivatives["original"].URL == nil {
					t.Error("expected URL to be set")
					return
				}
				if *photos[0].Derivatives["original"].URL != "https://example.com/checksum.jpg" {
					t.Errorf("URL = %q, want https://example.com/checksum.jpg (checksum should take priority)", *photos[0].Derivatives["original"].URL)
				}
			},
		},
		{
			name: "no matching URLs",
			photos: []Image{
				{
					PhotoGUID: "photo1",
					Derivatives: map[string]Derivative{
						"original": {Checksum: "unknown"},
					},
				},
			},
			allURLs: map[string]string{
				"other": "https://example.com/other.jpg",
			},
			verify: func(t *testing.T, photos []Image) {
				if photos[0].Derivatives["original"].URL != nil {
					t.Error("expected URL to be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EnrichPhotosWithURLs(tt.photos, tt.allURLs)
			tt.verify(t, tt.photos)
		})
	}
}
