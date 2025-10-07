// ABOUTME: Test suite for filename sanitization and download utilities
// ABOUTME: Validates cross-platform safe filename generation
package icloudalbum

import (
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(t *testing.T, result string)
	}{
		{
			name:  "remove special characters",
			input: "photo<>:\"/?*",
			validate: func(t *testing.T, result string) {
				if strings.ContainsAny(result, "<>:\"/?*") {
					t.Errorf("sanitize() = %q, should not contain special chars", result)
				}
			},
		},
		{
			name:  "normal filename unchanged",
			input: "my_photo_123",
			validate: func(t *testing.T, result string) {
				if result != "my_photo_123" {
					t.Errorf("sanitize() = %q, want my_photo_123", result)
				}
			},
		},
		{
			name:  "truncate long filenames",
			input: strings.Repeat("a", 250),
			validate: func(t *testing.T, result string) {
				if len(result) > 210 {
					t.Errorf("sanitize() length = %d, should be <= 210", len(result))
				}
				if !strings.HasSuffix(result, "_truncated") {
					t.Error("long filename should end with _truncated")
				}
			},
		},
		{
			name:  "trim leading/trailing dots and spaces",
			input: "  .photo.  ",
			validate: func(t *testing.T, result string) {
				if strings.HasPrefix(result, ".") || strings.HasPrefix(result, " ") {
					t.Errorf("sanitize() = %q, should not start with . or space", result)
				}
				if strings.HasSuffix(result, ".") || strings.HasSuffix(result, " ") {
					t.Errorf("sanitize() = %q, should not end with . or space", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitize(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestTrimDots(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" .photo. ", "photo"},
		{"..photo", "photo"},
		{"photo..", "photo"},
		{"photo", "photo"},
		{".....", ""},
		{"  ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := trimDots(tt.input)
			if result != tt.expected {
				t.Errorf("trimDots(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
