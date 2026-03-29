package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3.1 Safari/605.1.15"

var STATIC_PATHS []string = []string{
	"assets/css/global.css",
	"assets/css/dark.css",
	"assets/js/spin.js",
	"assets/menu-outline.svg",
	"assets/js/ionicons/index.esm.js",
	"assets/js/ionicons/ionicons.esm.js",
	"assets/js/ionicons/p-7a41fcdf.entry.js",
	"assets/js/ionicons/p-BKJPfAGl.js",
	"assets/js/ionicons/p-DQuL1Twl.js",
	"assets/js/ionicons/p-Z3yp5Yym.js",
	"assets/js/ionicons/svg/bulb-outline.svg",
	"assets/js/ionicons/svg/close-circle-outline.svg",
	"assets/js/ionicons/svg/menu-outline.svg",
	"assets/js/ionicons/svg/search-outline.svg",
	"assets/js/ionicons/svg/star-outline.svg",
	"assets/fonts/Lato-Light.ttf",
	"assets/apple-touch-icon.png",
	"favicon.png",
}

var postCnt = 0
var skipCnt = 0

var linkCnt = 0
var linkErrCnt = 0
var linkErrs = []error{}

var mutex = sync.Mutex{}

/*
Crawl site starting from specified page and save files to outDir,
optionally using password authentication.

If skipExisting, pages and assets already in outDir will be skipped.
If skipAssets, CSS, JS, and media assets will not be saved.
If checkLinks, check external links for availability without saving any files.
*/
func Crawl(page string, outDir string, password string, skipExisting bool,
	skipAssets bool, checkLinks bool) {
	const retries = 5
	attempts := sync.Map{}

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

	c := colly.NewCollector(colly.Async(true), colly.MaxBodySize(100*1024*1024))
	c.SetRequestTimeout(10 * time.Minute)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	private, err := IsPrivate()
	if err != nil {
		log.Fatal(err)
	}

	if private {
		if password == "" {
			flag.Usage()
			fmt.Println("")
			log.Fatal("Site is private and no password provided.")
		} else {
			err := Login(c, password)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	c.OnResponse(func(r *colly.Response) {
		if !checkLinks {
			SaveResponse(r)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		url := r.Request.URL.String()

		err := fmt.Errorf("Error fetching: %s: %d, %s", url, r.StatusCode, e)
		log.Print(err)

		if r.StatusCode != 200 {
			v, _ := attempts.LoadOrStore(url, 0)

			cnt := v.(int)
			if cnt < retries {
				attempts.Store(url, cnt+1)
				log.Printf("Retry %d/%d: %s", cnt+1, retries, url)
				r.Request.Retry()
			} else {
				log.Printf("%s: Max retries.", url)
			}
		}
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

	c.OnHTML("#cako-post-feed .cako-post-link", func(e *colly.HTMLElement) {
		mutex.Lock()
		postCnt++
		mutex.Unlock()

		if !Exists(e.Attr("href")) || !skipExisting {
			e.Request.Visit(e.Attr("href"))
		} else {
			mutex.Lock()
			skipCnt++
			mutex.Unlock()
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

		c.OnHTML("source", func(e *colly.HTMLElement) {
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

	if checkLinks {
		c.OnHTML("a", checkLink)
	}

	c.Visit(page)
	c.Wait()

	if checkLinks {
		fmt.Println()
		log.Printf("%d links visited, %d errored", linkCnt, linkErrCnt)
		for _, err := range linkErrs {
			fmt.Println(err)
		}
	} else {
		fmt.Println()
		log.Printf("%d posts found, %d saved", postCnt, postCnt-skipCnt)
	}

}
