# iCloud Album Go

A Go library and CLI tools for fetching photos from iCloud shared albums.

This is an idiomatic Go port of the Rust icloud-album crate, preserving all behavior including:
- Base URL calculation with base62 partition algorithm
- Apple's custom HTTP 330 redirect handling
- Schema-tolerant JSON parsing (handles both numeric and string values)
- Retry logic with exponential backoff and jitter
- URL enrichment (checksum and photo GUID matching)
- Best derivative selection
- MIME type detection with magic numbers
- Safe, cross-platform filename generation

## Installation

```bash
go get github.com/harperreed/icloud-album-go
```

## Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/harperreed/icloud-album-go/pkg/icloudalbum"
)

func main() {
    token := "your-album-token"

    resp, err := icloudalbum.GetICloudPhotos(token)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Album: %s\n", resp.Metadata.StreamName)
    fmt.Printf("Photos: %d\n", len(resp.Photos))
}
```

## Command-Line Tools

### album-info

Display basic album information:

```bash
go run ./cmd/album-info <shared_album_token>
```

### fetch-album

Fetch and display full album details with all photo derivatives:

```bash
go run ./cmd/fetch-album <shared_album_token>
```

### download-photos

Download all photos from an album:

```bash
go run ./cmd/download-photos <shared_album_token> <download_dir>
```

## Building

Build all command-line tools:

```bash
go build -o bin/album-info ./cmd/album-info
go build -o bin/fetch-album ./cmd/fetch-album
go build -o bin/download-photos ./cmd/download-photos
```

## Testing

Run all tests:

```bash
go test ./pkg/icloudalbum/... -v
```

## Project Structure

```
icloud-album-go/
  go.mod
  pkg/
    icloudalbum/
      models.go          # Data models with flexible JSON unmarshaling
      baseurl.go         # Base62 partition and URL calculation
      redirect.go        # Apple 330 redirect handling
      api.go             # API client with retry logic
      enrich.go          # Photo URL enrichment
      utils.go           # MIME detection and derivative selection
      download.go        # Photo download with filename sanitization
      icloud.go          # Main orchestrator
  cmd/
    album-info/main.go
    fetch-album/main.go
    download-photos/main.go
```

## Features

### Flexible JSON Parsing

Handles Apple's inconsistent API responses where fields can be either numbers or strings:

```go
type Uint64OrString uint64  // Accepts both 12345 and "12345"
type Uint32OrString uint32  // Accepts both 1920 and "1920"
```

### Retry Logic

Configurable retry with multiple backoff strategies:

- Constant delay
- Linear backoff
- Exponential backoff
- Exponential backoff with jitter (default)

### Smart Derivative Selection

Automatically selects the best quality photo:

1. Prefers originals ("original", "full", keys "3" or "4") with dimensions
2. Falls back to highest resolution available
3. Uses first available URL if no dimensions present

### Safe Filenames

Generates cross-platform safe filenames:

- Sanitizes special characters
- Truncates long names
- Trims leading/trailing dots and spaces
- Includes GUID, caption, and optional index

## Differences from Rust Version

- **Synchronous**: Uses blocking `net/http` (add contexts for cancellation if needed)
- **Lenient enrichment**: Handles both checksum and photo GUID URL keying
- **Same retry logic**: Exponential backoff with jitter, status-based retries
- **Same Apple 330 handling**: Literal status code 330 with JSON `X-Apple-MMe-Host` field

## License

MIT
