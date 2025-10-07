// ABOUTME: Test suite for iCloud base URL calculation using base62 partition algorithm
// ABOUTME: Verifies correct URL generation for different token formats
package icloudalbum

import (
	"testing"
)

func TestCharToBase62(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected uint32
		wantErr  bool
	}{
		{"digit 0", '0', 0, false},
		{"digit 5", '5', 5, false},
		{"digit 9", '9', 9, false},
		{"uppercase A", 'A', 10, false},
		{"uppercase Z", 'Z', 35, false},
		{"lowercase a", 'a', 36, false},
		{"lowercase z", 'z', 61, false},
		{"invalid char", '#', 0, true},
		{"invalid char", '@', 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := charToBase62(tt.char)
			if (err != nil) != tt.wantErr {
				t.Errorf("charToBase62(%c) error = %v, wantErr %v", tt.char, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("charToBase62(%c) = %v, want %v", tt.char, result, tt.expected)
			}
		})
	}
}

func TestCalculatePartition(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected uint32
		wantErr  bool
	}{
		{
			name:     "token starting with 0",
			token:    "0abc123",
			expected: 1, // (0 % 40) + 1 = 1
		},
		{
			name:     "token starting with A",
			token:    "Axyz789",
			expected: 11, // (10 % 40) + 1 = 11
		},
		{
			name:     "token starting with a",
			token:    "atoken123",
			expected: 37, // (36 % 40) + 1 = 37
		},
		{
			name:     "token starting with z",
			token:    "ztest",
			expected: 22, // (61 % 40) + 1 = 22
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid first char",
			token:   "#invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculatePartition(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculatePartition(%q) error = %v, wantErr %v", tt.token, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("calculatePartition(%q) = %v, want %v", tt.token, result, tt.expected)
			}
		})
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
		wantErr  bool
	}{
		{
			name:     "token starting with 0",
			token:    "0abc123",
			expected: "https://p01-sharedstreams.icloud.com/0abc123/sharedstreams/",
		},
		{
			name:     "token starting with A",
			token:    "Axyz789",
			expected: "https://p11-sharedstreams.icloud.com/Axyz789/sharedstreams/",
		},
		{
			name:     "token starting with a",
			token:    "atoken123",
			expected: "https://p37-sharedstreams.icloud.com/atoken123/sharedstreams/",
		},
		{
			name:     "token starting with 9",
			token:    "9test",
			expected: "https://p10-sharedstreams.icloud.com/9test/sharedstreams/",
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "#invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetBaseURL(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseURL(%q) error = %v, wantErr %v", tt.token, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetBaseURL(%q) = %q, want %q", tt.token, result, tt.expected)
			}
		})
	}
}

func TestGetBaseURL_ErrorConditions(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "empty token returns ErrEmptyToken",
			token:   "",
			wantErr: ErrEmptyToken,
		},
		{
			name:    "invalid char returns ErrInvalidBase62",
			token:   "!invalid",
			wantErr: ErrInvalidBase62,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBaseURL(tt.token)
			if err != tt.wantErr {
				t.Errorf("GetBaseURL(%q) error = %v, want %v", tt.token, err, tt.wantErr)
			}
		})
	}
}
