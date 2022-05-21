package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly/v2"
)

var STATIC_PATHS []string = []string{
	"assets/css/global.css",
	"assets/css/dark.css",
	"assets/js/spin.js",
	"assets/menu-outline.svg",
	"assets/js/ionicons/ionicons.esm.js",
	"assets/js/ionicons/ionicons.js",
	"assets/js/ionicons/p-3f680f7e.js",
	"assets/js/ionicons/p-5c60b45e.js",
	"assets/js/ionicons/p-e26ac56f.js",
	"assets/js/ionicons/svg/bulb-outline.svg",
	"assets/js/ionicons/svg/close-circle-outline.svg",
	"assets/js/ionicons/svg/menu-outline.svg",
	"assets/js/ionicons/svg/search-outline.svg",
	"assets/js/ionicons/svg/star-outline.svg",
	"assets/fonts/Lato-Light.ttf",
	"favicon.png",
	"apple-touch-icon.png",
}

// Check if local file exists based on full web URL
func Exists(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		panic("Error parsing url.")
	}

	path := strings.Trim(u.Path, "/")

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		return true
	}

	return false
}

func OnResponse(r *colly.Response) {
	var dest string

	if r.Request.URL.Path == "/" || r.Request.URL.Path == "/all/" {
		dest = "index.html"
	} else {
		dest = strings.Trim(r.Request.URL.Path, "/")
	}

	if dest != path.Base(dest) {
		destDir := path.Dir(dest)
		_, err := os.Stat(destDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(destDir, os.FileMode(0770))
			if err != nil {
				panic("Error creating output directory.")
			}
		}
	}

	r.Save(dest)
	fmt.Printf("Saved: %s\n", dest)
}

/* Crawl site starting from specified page and save files to outDir.

If skipExisting, pages and assets already in outDir will be skipped.
If skipAssets, CSS, JS, and media assets will not be saved.
*/
func Crawl(page string, outDir string, skipExisting bool, skipAssets bool) {
	postCnt := 0
	skipCnt := 0

	_, err := os.Stat(outDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(outDir, os.FileMode(0770))
		if err != nil {
			panic("Error creating output directory.")
		}
	}

	err = os.Chdir(outDir)
	if err != nil {
		log.Print(err)
		panic("Error changing to output directory.")
	}

	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		OnResponse(r)
	})

	// paths we need to manually visit
	if !skipAssets {
		for _, path := range STATIC_PATHS {
			uri := CAKO_IO_URL + path
			if !Exists(uri) || !skipExisting {
				c.Visit(uri)
			}
		}
	}

	if page == CAKO_IO_URL+"all/" {
		c.Visit(CAKO_IO_URL + "features/")
	}

	c.OnHTML(".cako-post-link", func(e *colly.HTMLElement) {
		postCnt++
		if !Exists(e.Attr("href")) || !skipExisting {
			e.Request.Visit(e.Attr("href"))
		} else {
			skipCnt++
		}
	})

	if !skipAssets {
		c.OnHTML("link[rel]", func(e *colly.HTMLElement) {
			rel := e.Attr("rel")
			if rel == "stylesheet" {
				if !Exists(e.Attr("href")) || !skipExisting {
					e.Request.Visit(e.Attr(("href")))
				}
			}
		})

		c.OnHTML("script", func(e *colly.HTMLElement) {
			if !Exists(e.Attr(("src"))) || !skipExisting {
				e.Request.Visit(e.Attr(("src")))
			}
		})

		c.OnHTML("img", func(e *colly.HTMLElement) {
			if !Exists(e.Attr(("src"))) || !skipExisting {
				e.Request.Visit(e.Attr(("src")))
			}
		})

		c.OnHTML("audio", func(e *colly.HTMLElement) {
			if !Exists(e.Attr(("src"))) || !skipExisting {
				e.Request.Visit(e.Attr(("src")))
			}
		})
	}

	c.Visit(page)

	summary := fmt.Sprintf("%d posts found, %d saved", postCnt, postCnt-skipCnt)

	fmt.Println(summary)
}
