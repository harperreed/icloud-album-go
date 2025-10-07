// ABOUTME: Test suite for iCloud album models, focusing on flexible JSON unmarshaling
// ABOUTME: Tests handle schema drift by accepting both numeric and string values
package icloudalbum

import (
	"encoding/json"
	"testing"
)

func TestUint64OrString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		wantErr  bool
	}{
		{
			name:     "numeric value",
			input:    `{"value": 12345}`,
			expected: 12345,
		},
		{
			name:     "string value",
			input:    `{"value": "67890"}`,
			expected: 67890,
		},
		{
			name:     "null value",
			input:    `{"value": null}`,
			expected: 0,
		},
		{
			name:     "empty string",
			input:    `{"value": ""}`,
			expected: 0,
		},
		{
			name:     "large number",
			input:    `{"value": 9223372036854775807}`,
			expected: 9223372036854775807,
		},
		{
			name:     "large string number",
			input:    `{"value": "9223372036854775807"}`,
			expected: 9223372036854775807,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				Value Uint64OrString `json:"value"`
			}
			err := json.Unmarshal([]byte(tt.input), &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if uint64(result.Value) != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", uint64(result.Value), tt.expected)
			}
		})
	}
}

func TestUint32OrString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint32
		wantErr  bool
	}{
		{
			name:     "numeric value",
			input:    `{"value": 12345}`,
			expected: 12345,
		},
		{
			name:     "string value",
			input:    `{"value": "67890"}`,
			expected: 67890,
		},
		{
			name:     "null value",
			input:    `{"value": null}`,
			expected: 0,
		},
		{
			name:     "empty string",
			input:    `{"value": ""}`,
			expected: 0,
		},
		{
			name:     "max uint32",
			input:    `{"value": 4294967295}`,
			expected: 4294967295,
		},
		{
			name:     "max uint32 as string",
			input:    `{"value": "4294967295"}`,
			expected: 4294967295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				Value Uint32OrString `json:"value"`
			}
			err := json.Unmarshal([]byte(tt.input), &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if uint32(result.Value) != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", uint32(result.Value), tt.expected)
			}
		})
	}
}

func TestDerivative_UnmarshalJSON(t *testing.T) {
	input := `{
		"checksum": "abc123",
		"fileSize": "1024",
		"width": 1920,
		"height": "1080",
		"url": "https://example.com/photo.jpg"
	}`

	var d Derivative
	err := json.Unmarshal([]byte(input), &d)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if d.Checksum != "abc123" {
		t.Errorf("Checksum = %v, want abc123", d.Checksum)
	}
	if d.FileSize == nil || uint64(*d.FileSize) != 1024 {
		t.Errorf("FileSize = %v, want 1024", d.FileSize)
	}
	if d.Width == nil || uint32(*d.Width) != 1920 {
		t.Errorf("Width = %v, want 1920", d.Width)
	}
	if d.Height == nil || uint32(*d.Height) != 1080 {
		t.Errorf("Height = %v, want 1080", d.Height)
	}
	if d.URL == nil || *d.URL != "https://example.com/photo.jpg" {
		t.Errorf("URL = %v, want https://example.com/photo.jpg", d.URL)
	}
}

func TestImage_UnmarshalJSON(t *testing.T) {
	input := `{
		"photoGuid": "photo-123",
		"derivatives": {
			"original": {
				"checksum": "abc123",
				"fileSize": 2048,
				"width": "3840",
				"height": 2160
			}
		},
		"caption": "Test Photo",
		"dateCreated": "2024-01-01",
		"width": 3840,
		"height": "2160"
	}`

	var img Image
	err := json.Unmarshal([]byte(input), &img)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if img.PhotoGUID != "photo-123" {
		t.Errorf("PhotoGUID = %v, want photo-123", img.PhotoGUID)
	}
	if len(img.Derivatives) != 1 {
		t.Errorf("len(Derivatives) = %v, want 1", len(img.Derivatives))
	}
	if img.Caption == nil || *img.Caption != "Test Photo" {
		t.Errorf("Caption = %v, want Test Photo", img.Caption)
	}
	if img.Width == nil || uint32(*img.Width) != 3840 {
		t.Errorf("Width = %v, want 3840", img.Width)
	}
	if img.Height == nil || uint32(*img.Height) != 2160 {
		t.Errorf("Height = %v, want 2160", img.Height)
	}
}

func TestApiResponse_UnmarshalJSON(t *testing.T) {
	input := `{
		"photos": [{
			"photoGuid": "photo-1",
			"derivatives": {}
		}],
		"photoGuids": ["photo-1"],
		"streamName": "Test Album",
		"userFirstName": "John",
		"userLastName": "Doe",
		"streamCtag": "tag123",
		"itemsReturned": "1"
	}`

	var resp ApiResponse
	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if len(resp.Photos) != 1 {
		t.Errorf("len(Photos) = %v, want 1", len(resp.Photos))
	}
	if len(resp.PhotoGuids) != 1 {
		t.Errorf("len(PhotoGuids) = %v, want 1", len(resp.PhotoGuids))
	}
	if resp.StreamName == nil || *resp.StreamName != "Test Album" {
		t.Errorf("StreamName = %v, want Test Album", resp.StreamName)
	}
	if resp.UserFirstName == nil || *resp.UserFirstName != "John" {
		t.Errorf("UserFirstName = %v, want John", resp.UserFirstName)
	}
	if resp.ItemsReturned == nil || uint32(*resp.ItemsReturned) != 1 {
		t.Errorf("ItemsReturned = %v, want 1", resp.ItemsReturned)
	}
}
