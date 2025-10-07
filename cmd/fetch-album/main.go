// ABOUTME: Command-line tool to fetch and display full iCloud album details
// ABOUTME: Shows all photos with their derivatives, dimensions, and URLs
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harperreed/icloud-album-go/pkg/icloudalbum"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: fetch-album <shared_album_token>")
		os.Exit(2)
	}
	token := os.Args[1]

	resp, err := icloudalbum.GetICloudPhotos(token)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("\nAlbum: %s\n", resp.Metadata.StreamName)
	fmt.Printf("Owner: %s %s\n", resp.Metadata.UserFirstName, resp.Metadata.UserLastName)
	fmt.Printf("Photos: %d\n", len(resp.Photos))
	for i, p := range resp.Photos {
		fmt.Printf("\nPhoto %d: %s\n", i+1, p.PhotoGUID)
		for k, d := range p.Derivatives {
			w, h := uint32(0), uint32(0)
			if d.Width != nil {
				w = uint32(*d.Width)
			}
			if d.Height != nil {
				h = uint32(*d.Height)
			}
			u := "No URL"
			if d.URL != nil {
				u = *d.URL
			}
			fmt.Printf("  %s: %dx%d  %s\n", k, w, h, u)
		}
	}
}
