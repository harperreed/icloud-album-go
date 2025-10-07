// ABOUTME: Handles Apple's custom HTTP 330 redirect status for iCloud shared streams
// ABOUTME: Extracts redirected host from JSON response to build correct base URL
package icloudalbum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRedirectedBaseURL detects Apple's custom 330 redirect and, if present,
// constructs https://{host}/{token}/sharedstreams/. Otherwise returns baseURL.
func GetRedirectedBaseURL(client *http.Client, baseURL, token string) (string, error) {
	type payload struct {
		StreamCTag *string `json:"streamCtag"` // null
	}
	body, _ := json.Marshal(payload{StreamCTag: nil})
	req, err := http.NewRequest("POST", baseURL+"webstream", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Apple's server uses status 330 (non-standard) to signal redirect host in JSON.
	if resp.StatusCode == 330 {
		var m map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return "", err
		}
		if host, _ := m["X-Apple-MMe-Host"].(string); host != "" {
			return fmt.Sprintf("https://%s/%s/sharedstreams/", host, token), nil
		}
	}
	// No redirect or missing host → use original
	return baseURL, nil
}
