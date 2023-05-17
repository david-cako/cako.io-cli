package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const CAKO_IO_URL = "https://cako.io/"

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
	serve := flag.Bool("serve", false, "Serve locally saved files")

	flag.Parse()

	if *serve {
		fmt.Printf("serving on http://localhost:8080\n")
		http.Handle("/", http.FileServer(http.Dir(*outDir)))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	Crawl(CAKO_IO_URL+*page, *outDir, *password, *skipExisting, *skipAssets)
}
