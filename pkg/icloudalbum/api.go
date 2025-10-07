// ABOUTME: API client for iCloud shared streams with retry logic and exponential backoff
// ABOUTME: Fetches album metadata, photos, and asset URLs with configurable retry strategies
package icloudalbum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

// getAPIResponse performs POST {base}/webstream and returns parsed Photos + Metadata.
func getAPIResponse(client *http.Client, baseURL string) ([]Image, Metadata, error) {
	type payload struct {
		StreamCTag *string `json:"streamCtag"` // null
	}
	body, _ := json.Marshal(payload{StreamCTag: nil})

	req, err := http.NewRequest("POST", baseURL+"webstream", bytes.NewReader(body))
	if err != nil {
		return nil, Metadata{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, Metadata{}, fmt.Errorf("webstream request failed (status %d)", resp.StatusCode)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, Metadata{}, err
	}

	// Lenient parse into ApiResponse
	var api ApiResponse
	if err := json.Unmarshal(raw, &api); err != nil {
		log.Printf("warn: error deserializing API response: %v", err)
	}

	// streamName is required for a valid album (mirror Rust's Required severity)
	if api.StreamName == nil || *api.StreamName == "" {
		return nil, Metadata{}, errors.New("missing required field: streamName")
	}

	// Build metadata (fallbacks mirror the Rust behavior)
	items := uint32(0)
	if api.ItemsReturned != nil {
		items = uint32(*api.ItemsReturned)
	}
	locs := json.RawMessage("null")
	if api.Locations != nil {
		locs = *api.Locations
	}

	md := Metadata{
		StreamName:    *api.StreamName,
		UserFirstName: derefOr(api.UserFirstName, ""),
		UserLastName:  derefOr(api.UserLastName, ""),
		StreamCTag:    derefOr(api.StreamCTag, ""),
		ItemsReturned: items,
		Locations:     locs,
	}

	return api.Photos, md, nil
}

func derefOr[T ~string](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// ------------------------------ Asset URLs + Retry ----------------------------

type BackoffStrategy int

const (
	BackoffConstant BackoffStrategy = iota
	BackoffLinear
	BackoffExponential
	BackoffExponentialWithJitter
)

type RetryConfig struct {
	MaxRetries                 int
	BaseDelay                  time.Duration
	Strategy                   BackoffStrategy
	MaxDelay                   time.Duration
	RetryableStatusCodes       []int // specific codes
	PermanentFailureStatusCodes []int
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  500 * time.Millisecond,
		Strategy:   BackoffExponentialWithJitter,
		MaxDelay:   30 * time.Second,
		RetryableStatusCodes: []int{408, 429, 500, 502, 503, 504},
		PermanentFailureStatusCodes: []int{400, 401, 403, 404},
	}
}

func shouldRetryStatus(cfg RetryConfig, code int) bool {
	for _, s := range cfg.PermanentFailureStatusCodes {
		if code == s {
			return false
		}
	}
	for _, s := range cfg.RetryableStatusCodes {
		if code == s {
			return true
		}
	}
	return code >= 500 && code <= 599
}

func nextDelay(cfg RetryConfig, attempt int) time.Duration {
	switch cfg.Strategy {
	case BackoffConstant:
		return cfg.BaseDelay
	case BackoffLinear:
		d := time.Duration(attempt) * cfg.BaseDelay
		if d > cfg.MaxDelay {
			return cfg.MaxDelay
		}
		return d
	case BackoffExponential:
		d := cfg.BaseDelay * (1 << min(attempt, 30))
		if d > cfg.MaxDelay {
			return cfg.MaxDelay
		}
		return d
	case BackoffExponentialWithJitter:
		max := cfg.BaseDelay * (1 << min(attempt, 30))
		if max > cfg.MaxDelay {
			max = cfg.MaxDelay
		}
		if max <= 0 {
			return 0
		}
		return time.Duration(rand.Int64N(max.Milliseconds()+1)) * time.Millisecond
	default:
		return cfg.BaseDelay
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetAssetURLs calls {base}/webasseturls with photo GUIDs and returns a map id->fullURL.
// Note: "id" keys are whatever Apple returns in `items` (photoGuid or checksum).
func GetAssetURLs(client *http.Client, baseURL string, photoGUIDs []string, cfg *RetryConfig) (map[string]string, error) {
	c := DefaultRetryConfig()
	if cfg != nil {
		c = *cfg
	}
	if len(photoGUIDs) == 0 {
		log.Printf("warn: get_asset_urls called with empty photoGUIDs")
		return map[string]string{}, nil
	}

	type payload struct {
		PhotoGuids []string `json:"photoGuids"`
	}
	body, _ := json.Marshal(payload{PhotoGuids: photoGUIDs})

	attempt := 0
	for {
		req, err := http.NewRequest("POST", baseURL+"webasseturls", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			// treat as retryable network error
			if attempt >= c.MaxRetries {
				return nil, fmt.Errorf("webasseturls network error after retries: %w", err)
			}
			time.Sleep(nextDelay(c, attempt))
			attempt++
			continue
		}

		// Special handling: 400 → known Apple quirk; continue with empty map (parity with Rust)
		if resp.StatusCode == 400 {
			resp.Body.Close()
			log.Printf("warn: webasseturls returned 400; returning empty map for partial functionality")
			return map[string]string{}, nil
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			if shouldRetryStatus(c, resp.StatusCode) && attempt < c.MaxRetries {
				time.Sleep(nextDelay(c, attempt))
				attempt++
				continue
			}
			return nil, fmt.Errorf("webasseturls request failed (status %d)", resp.StatusCode)
		}

		// Parse successful response
		var parsed struct {
			Items map[string]struct {
				URLLocation string `json:"url_location"`
				URLPath     string `json:"url_path"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		res := make(map[string]string, len(parsed.Items))
		for id, it := range parsed.Items {
			if it.URLLocation == "" || it.URLPath == "" {
				log.Printf("warn: missing url_location or url_path for id %s", id)
				continue
			}
			res[id] = "https://" + it.URLLocation + it.URLPath
		}
		return res, nil
	}
}
