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
