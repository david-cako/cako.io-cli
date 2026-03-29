package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Returns true if home page redirects to login.
func IsPrivate() (bool, error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	r, err := client.Get(CAKO_IO_URL)
	if err != nil {
		return false, err
	}

	loc, err := r.Location()
	if err != nil {
		return false, err
	}

	if r.StatusCode == http.StatusFound && strings.Contains(loc.Path, "private") {
		return true, nil
	}

	return false, nil
}

// Lol
func Login(collector *colly.Collector, password string) error {
	loginUrl := CAKO_IO_URL + "private/"
	ghostCookieName := "ghost-private"

	client := http.Client{
		/* Keep 302 response so we can extract express session cookie */
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	params := url.Values{}
	params.Add("password", password)

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	setCookie := r.Header.Get("Set-Cookie")

	if !strings.Contains(setCookie, ghostCookieName) {
		return fmt.Errorf("Login failed.")
	}

	collector.SetCookies(CAKO_IO_URL, r.Cookies())

	return nil
}
