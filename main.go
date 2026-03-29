package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const CAKO_IO_URL = "https://cako.io/"

var YOUTUBE_DATA_API_KEY string
var ytService *youtube.Service

func main() {
	fmt.Printf("\ncako.io cli\n\n")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command line interface for crawling "+
			"and locally serving pages from cako.io\n\n")

		fmt.Fprintf(os.Stderr, "Usage: \n")
		flag.PrintDefaults()
	}
	page := flag.String("page", "all/", "Crawl specified page name only")
	outDir := flag.String("outDir", "./saved/", "Output directory to save files")
	password := flag.String("password", "", "Password for private site access")
	skipExisting := flag.Bool("skipExisting", false, "Skip crawling pages already in output directory")
	skipAssets := flag.Bool("skipAssets", false, "Only crawl html files")
	checkLinks := flag.Bool("checkLinks", false, "Only check links for reachability (do not archive files).")
	ytDataApiKey := flag.String("youtubeDataApiKey", "",
		"YouTube Data API key for verifying YouTube video availability.")
	serve := flag.Bool("serve", false, "Serve locally saved files")

	flag.Parse()

	if *serve {
		fmt.Printf("serving on http://localhost:8080\n")
		http.Handle("/", http.FileServer(http.Dir(*outDir)))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	if *ytDataApiKey != "" {
		YOUTUBE_DATA_API_KEY = *ytDataApiKey
	} else {
		YOUTUBE_DATA_API_KEY = os.Getenv("YOUTUBE_DATA_API_KEY")
	}

	if YOUTUBE_DATA_API_KEY != "" {
		s, err := youtube.NewService(context.Background(),
			option.WithAPIKey(YOUTUBE_DATA_API_KEY))
		if err != nil {
			log.Fatalf("Could not create YouTube client: %v", err)
		}

		ytService = s
	}

	Crawl(CAKO_IO_URL+*page, *outDir, *password, *skipExisting, *skipAssets, *checkLinks)
}
