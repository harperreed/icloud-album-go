// ABOUTME: Generates iCloud shared streams base URLs using base62 partition calculation
// ABOUTME: Maps album tokens to correct regional server endpoints (p01-p40)
package icloudalbum

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyToken      = errors.New("empty token")
	ErrInvalidBase62   = errors.New("invalid base62 character")
)

func charToBase62(r rune) (uint32, error) {
	switch {
	case r >= '0' && r <= '9':
		return uint32(r - '0'), nil
	case r >= 'A' && r <= 'Z':
		return uint32(r-'A') + 10, nil
	case r >= 'a' && r <= 'z':
		return uint32(r-'a') + 36, nil
	default:
		return 0, ErrInvalidBase62
	}
}

func calculatePartition(token string) (uint32, error) {
	if token == "" || strings.TrimSpace(token) == "" {
		return 0, ErrEmptyToken
	}
	r := []rune(token)[0]
	v, err := charToBase62(r)
	if err != nil {
		return 0, err
	}
	return 1 + (v % 40), nil
}

// GetBaseURL builds: https://pXX-sharedstreams.icloud.com/{token}/sharedstreams/
func GetBaseURL(token string) (string, error) {
	part, err := calculatePartition(token)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://p%02d-sharedstreams.icloud.com/%s/sharedstreams/", part, token), nil
}
