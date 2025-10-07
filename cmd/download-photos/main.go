// ABOUTME: Command-line tool to download all photos from an iCloud shared album
// ABOUTME: Saves photos with descriptive filenames including GUIDs and captions
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/harperreed/icloud-album-go/pkg/icloudalbum"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: download-photos <shared_album_token> <download_dir>")
		os.Exit(2)
	}
	token := os.Args[1]
	outDir := os.Args[2]

	resp, err := icloudalbum.GetICloudPhotos(token)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Album: %s (%d photos)\n", resp.Metadata.StreamName, len(resp.Photos))

	for i := range resp.Photos {
		p := &resp.Photos[i]
		fmt.Printf("Downloading %d/%d: %s\n", i+1, len(resp.Photos), p.PhotoGUID)
		fp, err := icloudalbum.DownloadPhoto(p, &i, outDir, nil)
		if err != nil {
			fmt.Printf("  failed: %v\n", err)
			continue
		}
		fmt.Printf("  saved: %s\n", filepath.Base(fp))
	}
}
