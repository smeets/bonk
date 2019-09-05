package main

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
)

func getevent(url string) (*goquery.Document, error) {
	www := "https://www.ufc.com"

	if !strings.HasPrefix(url, "https://") {
		url = www + url
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", "UFCCookieTest=1; STYXKEY_region=SWEDEN.en-se.Europe/Stockholm")
	req.Header.Set("Accept-Language", "en-se")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func findevents() []string {
	www := "https://www.ufc.com/events"
	// 
	req, err := http.NewRequest("GET", www, nil)
	bailout("newreq:", err)
	req.Header.Set("Cookie", "UFCCookieTest=1; STYXKEY_region=SWEDEN.en-se.Europe/Stockholm")
	req.Header.Set("Accept-Language", "en-se")
	res, err := http.DefaultClient.Do(req)
	bailout("findevents:", err)

	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	bailout("newdoc:", err)

	eventlink := "#events-list-upcoming .c-card-event--result__info a"
	
	var links []string
	doc.Find(eventlink).Each(func(n int, sel *goquery.Selection) {
		link, ok := sel.Attr("href")
		if !ok {
			return
		}
		links = append(links, link)
	})

	return links
}