// ABOUTME: Command-line tool to display basic iCloud shared album information
// ABOUTME: Shows album name, owner, photo count, and first few photos
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
		fmt.Fprintln(os.Stderr, "usage: album-info <shared_album_token>")
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

	if len(resp.Photos) > 0 {
		fmt.Println("\nFirst few photos:")
		limit := 5
		if len(resp.Photos) < limit {
			limit = len(resp.Photos)
		}
		for i := 0; i < limit; i++ {
			p := resp.Photos[i]
			date := "N/A"
			if p.DateCreated != nil {
				date = *p.DateCreated
			}
			caption := "N/A"
			if p.Caption != nil {
				caption = *p.Caption
			}
			fmt.Printf("  %2d  %s  %s\n", i+1, date, caption)
		}
	}
}
