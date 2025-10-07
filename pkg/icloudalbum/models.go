// ABOUTME: Data models for iCloud shared album API responses with flexible JSON unmarshaling
// ABOUTME: Handles schema drift by accepting both numeric and string representations
package icloudalbum

import (
	"encoding/json"
	"log"
	"strconv"
)

// -- Number-or-string helpers --------------------------------------------------

// Uint64OrString decodes a JSON number OR a quoted number into uint64.
type Uint64OrString uint64

func (u *Uint64OrString) UnmarshalJSON(b []byte) error {
	// null → leave pointer nil (caller uses *Uint64OrString)
	if string(b) == "null" {
		return nil
	}
	// Try as number
	var n uint64
	if err := json.Unmarshal(b, &n); err == nil {
		*u = Uint64OrString(n)
		return nil
	}
	// Try as string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			return nil
		}
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			log.Printf("warn: failed to parse string %q as u64: %v", s, err)
			return nil
		}
		*u = Uint64OrString(v)
		return nil
	}
	// Unknown type → ignore (schema drift)
	return nil
}

type Uint32OrString uint32

func (u *Uint32OrString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	// Try as number
	var n uint64
	if err := json.Unmarshal(b, &n); err == nil {
		*u = Uint32OrString(uint32(n))
		return nil
	}
	// Try as string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			return nil
		}
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			log.Printf("warn: failed to parse string %q as u32: %v", s, err)
			return nil
		}
		*u = Uint32OrString(uint32(v))
		return nil
	}
	return nil
}

// -- Models --------------------------------------------------------------------

type Derivative struct {
	Checksum string           `json:"checksum"`
	FileSize *Uint64OrString  `json:"fileSize,omitempty"`
	Width    *Uint32OrString  `json:"width,omitempty"`
	Height   *Uint32OrString  `json:"height,omitempty"`
	URL      *string          `json:"url,omitempty"`
}

type Image struct {
	PhotoGUID        string                     `json:"photoGuid"`
	Derivatives      map[string]Derivative      `json:"derivatives"`
	Caption          *string                    `json:"caption,omitempty"`
	DateCreated      *string                    `json:"dateCreated,omitempty"`
	BatchDateCreated *string                    `json:"batchDateCreated,omitempty"`
	Width            *Uint32OrString            `json:"width,omitempty"`
	Height           *Uint32OrString            `json:"height,omitempty"`
}

type Metadata struct {
	StreamName     string          `json:"streamName"`
	UserFirstName  string          `json:"userFirstName"`
	UserLastName   string          `json:"userLastName"`
	StreamCTag     string          `json:"streamCtag"`
	ItemsReturned  uint32          `json:"itemsReturned"`
	Locations      json.RawMessage `json:"locations"`
}

type ApiResponse struct {
	Photos        []Image          `json:"photos"`
	PhotoGuids    []string         `json:"photoGuids"`
	StreamName    *string          `json:"streamName,omitempty"`
	UserFirstName *string          `json:"userFirstName,omitempty"`
	UserLastName  *string          `json:"userLastName,omitempty"`
	StreamCTag    *string          `json:"streamCtag,omitempty"`
	ItemsReturned *Uint32OrString  `json:"itemsReturned,omitempty"`
	Locations     *json.RawMessage `json:"locations,omitempty"`
}

type ICloudResponse struct {
	Metadata Metadata
	Photos   []Image
}
