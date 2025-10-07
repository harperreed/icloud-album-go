// ABOUTME: Test suite for MIME detection and derivative selection utilities
// ABOUTME: Validates magic number detection and best quality derivative selection
package icloudalbum

import (
	"testing"
)

func TestExtensionFromMIME(t *testing.T) {
	tests := []struct {
		mime     string
		expected string
	}{
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"image/heic", ".heic"},
		{"image/heif", ".heif"},
		{"video/mp4", ".mp4"},
		{"video/quicktime", ".mov"},
		{"image/gif", ".gif"},
		{"unknown/type", ".jpg"}, // defaults to .jpg
	}

	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			result := ExtensionFromMIME(tt.mime)
			if result != tt.expected {
				t.Errorf("ExtensionFromMIME(%q) = %q, want %q", tt.mime, result, tt.expected)
			}
		})
	}
}

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		filename string
		expected string
	}{
		{
			name:     "JPEG magic number",
			bytes:    []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00},
			expected: "image/jpeg",
		},
		{
			name:     "PNG magic number",
			bytes:    []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: "image/png",
		},
		{
			name:     "GIF magic number",
			bytes:    []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61},
			expected: "image/gif",
		},
		{
			name:     "MP4 magic number",
			bytes:    []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x69, 0x73, 0x6F, 0x6D},
			expected: "video/mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectMIMEType(tt.bytes, tt.filename)
			if result != tt.expected {
				t.Errorf("DetectMIMEType() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSelectBestDerivative(t *testing.T) {
	w1920 := Uint32OrString(1920)
	w3840 := Uint32OrString(3840)
	h1080 := Uint32OrString(1080)
	h2160 := Uint32OrString(2160)
	url1 := "https://example.com/photo1.jpg"
	url2 := "https://example.com/photo2.jpg"

	tests := []struct {
		name        string
		derivatives map[string]Derivative
		wantKey     string
		wantURL     string
		wantOK      bool
	}{
		{
			name:        "empty derivatives",
			derivatives: map[string]Derivative{},
			wantOK:      false,
		},
		{
			name: "prefer original with dimensions",
			derivatives: map[string]Derivative{
				"thumbnail": {Checksum: "thumb", Width: &w1920, Height: &h1080, URL: &url1},
				"original":  {Checksum: "orig", Width: &w3840, Height: &h2160, URL: &url2},
			},
			wantKey: "original",
			wantURL: url2,
			wantOK:  true,
		},
		{
			name: "prefer full quality",
			derivatives: map[string]Derivative{
				"thumbnail": {Checksum: "thumb", Width: &w1920, Height: &h1080, URL: &url1},
				"full":      {Checksum: "full", Width: &w3840, Height: &h2160, URL: &url2},
			},
			wantKey: "full",
			wantURL: url2,
			wantOK:  true,
		},
		{
			name: "prefer key 3 or 4",
			derivatives: map[string]Derivative{
				"1": {Checksum: "1", Width: &w1920, Height: &h1080, URL: &url1},
				"4": {Checksum: "4", Width: &w3840, Height: &h2160, URL: &url2},
			},
			wantKey: "4",
			wantURL: url2,
			wantOK:  true,
		},
		{
			name: "no URLs available",
			derivatives: map[string]Derivative{
				"original": {Checksum: "orig", Width: &w3840, Height: &h2160},
			},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, _, url, ok := SelectBestDerivative(tt.derivatives)
			if ok != tt.wantOK {
				t.Errorf("SelectBestDerivative() ok = %v, want %v", ok, tt.wantOK)
				return
			}
			if !tt.wantOK {
				return
			}
			if key != tt.wantKey {
				t.Errorf("SelectBestDerivative() key = %q, want %q", key, tt.wantKey)
			}
			if url != tt.wantURL {
				t.Errorf("SelectBestDerivative() url = %q, want %q", url, tt.wantURL)
			}
		})
	}
}
