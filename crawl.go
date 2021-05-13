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

	if r.Request.URL.Path == "/" {
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
	pageCnt := 0
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

	c.OnHTML(".cako-post-link", func(e *colly.HTMLElement) {
		pageCnt++
		if !skipExisting || !Exists(e.Attr("href")) {
			e.Request.Visit(e.Attr("href"))
		} else {
			skipCnt++
		}
	})

	if !skipAssets {
		c.OnHTML("link[rel]", func(e *colly.HTMLElement) {
			rel := e.Attr("rel")
			if rel == "stylesheet" {
				if !skipExisting || !Exists(e.Attr("href")) {
					e.Request.Visit(e.Attr(("href")))
				}
			}
		})

		c.OnHTML("script", func(e *colly.HTMLElement) {
			if !skipExisting || !Exists(e.Attr(("src"))) {
				e.Request.Visit(e.Attr(("src")))
			}
		})

		c.OnHTML("img", func(e *colly.HTMLElement) {
			if !skipExisting || !Exists(e.Attr(("src"))) {
				e.Request.Visit(e.Attr(("src")))
			}
		})

		c.OnHTML("audio", func(e *colly.HTMLElement) {
			if !skipExisting || !Exists(e.Attr(("src"))) {
				e.Request.Visit(e.Attr(("src")))
			}
		})
	}

	c.OnResponse(func(r *colly.Response) {
		OnResponse(r)
	})

	// imports we need to manually visit
	if !skipAssets {
		globalCss := CAKO_IO_URL + "assets/css/global.css"
		if !skipExisting || !Exists(globalCss) {
			c.Visit(globalCss)
		}

		darkCss := CAKO_IO_URL + "assets/css/dark.css"
		if !skipExisting || !Exists(darkCss) {
			c.Visit(darkCss)
		}

		spinJs := CAKO_IO_URL + "assets/js/spin.js"
		if !skipExisting || !Exists(spinJs) {
			c.Visit(spinJs)
		}

		favicon := CAKO_IO_URL + "favicon.png"
		if !skipExisting || !Exists(favicon) {
			c.Visit(favicon)
		}

		appleTouchIcon := CAKO_IO_URL + "apple-touch-icon.png"
		if !skipExisting || !Exists(appleTouchIcon) {
			c.Visit(appleTouchIcon)
		}
	}

	c.Visit(page)

	summary := fmt.Sprintf("%d pages found", pageCnt)

	if skipExisting {
		summary += fmt.Sprintf(", %d skipped", skipCnt)
	}

	fmt.Println(summary)
}
