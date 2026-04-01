package main

import (
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/gocolly/colly/v2"
)

var STATIC_PATHS []string = []string{
	"assets/css/global.css",
	"assets/css/dark.css",
	"assets/css/private.css",
	"assets/menu-outline.svg",
	"assets/js/Api.js",
	"assets/js/CakoApp.js",
	"assets/js/Header.js",
	"assets/js/Html.js",
	"assets/js/InfiniteScroll.js",
	"assets/js/Lights.js",
	"assets/js/Menu.js",
	"assets/js/Search.js",
	"assets/js/lib/ionicons/index.esm.js",
	"assets/js/lib/ionicons/ionicons.esm.js",
	"assets/js/lib/ionicons/p-7a41fcdf.entry.js",
	"assets/js/lib/ionicons/p-BKJPfAGl.js",
	"assets/js/lib/ionicons/p-DQuL1Twl.js",
	"assets/js/lib/ionicons/p-Z3yp5Yym.js",
	"assets/js/lib/ionicons/svg/bulb-outline.svg",
	"assets/js/lib/ionicons/svg/close-circle-outline.svg",
	"assets/js/lib/ionicons/svg/menu-outline.svg",
	"assets/js/lib/ionicons/svg/search-outline.svg",
	"assets/js/lib/ionicons/svg/star-outline.svg",
	"assets/js/lib/content-api.min.js.map",
	"assets/fonts/Lato-Light.ttf",
	"assets/apple-touch-icon.png",
	"assets/favicon.ico",
	"assets/cako_rounded-70.png",
	"assets/cako_rounded.png",
	"assets/site.webmanifest",
}

func SaveResponse(r *colly.Response) {
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

	body := r.Body

	if utf8.Valid(r.Body) {
		body = ReplacePageTitle(r.Body)
		body = ReplaceContentUrls(r.Body)
	}

	os.WriteFile(dest, body, 0644)

	log.Printf("Saved: %s\n", dest)
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

func ReplacePageTitle(postBody []byte) []byte {
	s := string(postBody)
	s = strings.Replace(s, "<title>cako.io (Page 1)</title>", "<title>cako.io</title>", 1)

	return []byte(s)
}

func ReplaceContentUrls(postBody []byte) []byte {
	s := string(postBody)

	s = strings.ReplaceAll(s, "https://cako.io/content", "/content")
	s = strings.ReplaceAll(s, "http://cako.io/content", "/content")

	return []byte(s)
}
