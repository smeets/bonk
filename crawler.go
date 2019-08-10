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

	res, err := http.Get(url)
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

	res, err := http.Get(www)
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