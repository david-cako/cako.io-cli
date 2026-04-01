package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

func verifyPageAvailability(u url.URL) error {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
	}

	return nil
}

func verifyYoutubeAvailability(u url.URL) error {
	if !strings.Contains(u.Host, "youtube.com") {
		return errors.New("Not a YouTube URL.")
	}

	id := u.Query().Get("v")
	if id == "" {
		return errors.New("Missing video id in URL.")
	}

	if ytService == nil {
		return errors.New("YouTube service not created.  Set YOUTUBE_DATA_API_KEY variable.")
	}

	l := ytService.Videos.List([]string{"status"})
	l.Id(id)

	r, err := l.Do()
	if err != nil {
		return fmt.Errorf("Error calling youtube.videos.list: %w", err)
	}

	if len(r.Items) < 1 {
		return errors.New("Video not found.")
	}

	v := r.Items[0]
	uploadStatus := v.Status.UploadStatus
	privacyStatus := v.Status.PrivacyStatus

	if uploadStatus != "processed" && uploadStatus != "uploaded" {
		return fmt.Errorf("Video upload status: %s", uploadStatus)
	}

	if privacyStatus == "private" {
		return fmt.Errorf("Video is private.")
	}

	return nil
}

func checkLink(e *colly.HTMLElement) {
	u, err := url.Parse(e.Attr("href"))
	if err != nil {
		log.Printf("Error parsing link: %s, %s", e.Attr("href"), err)
		return
	}

	var altText string

	altText, _ = e.DOM.Find("img").Attr("alt")

	var visited bool
	var linkErr error

	if strings.Contains(u.Host, "youtube.com") {
		visited = true

		err := verifyYoutubeAvailability(*u)
		if err != nil {
			linkErr = err
		}

	} else if u.Host != "cako.io" && u.Host != "" && u.Host != "davidcako.com" {
		visited = true

		err := verifyPageAvailability(*u)
		if err != nil {
			linkErr = err
		}
	}

	if visited {
		mutex.Lock()
		linkCnt++
		mutex.Unlock()

		if altText != "" {
			altText = fmt.Sprintf(" (%s)", altText)
		}

		if linkErr != nil {
			le := fmt.Errorf("Error: %s: %s%s: %s", e.Request.URL.String(), u.String(), altText, linkErr)
			log.Println(le)

			mutex.Lock()
			linkErrs = append(linkErrs, le)
			mutex.Unlock()
		} else {
			log.Printf("%s: %s%s: OK", e.Request.URL.String(), u.String(), altText)
		}
	}
}
