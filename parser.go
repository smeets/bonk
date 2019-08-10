package main

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Event struct {
	Title, Headline, Venue, Start string
	URL  string
	Cards []Card
	Fights uint
}

type Card struct {
	root 	*goquery.Selection
	Name 	string
	Fights 	[]Fight
}

type Fight struct {
	Blue, Red Corner
}

type Corner struct {
	Name, Image string
}

func getpict(sel *goquery.Selection, corner string) string {
	class := ".c-listing-fight__corner-image--" + corner
	img := sel.Find(class).Find("img")
	src, ok := img.Attr("src")
	if !ok {
		return "n/a"
	}
	return src
}

func getname(sel *goquery.Selection, corner string) string {
	class := ".c-listing-fight__corner-body--" + corner
	body := sel.Find(class)
	given := body.Find(".c-listing-fight__corner-given-name").Text()
	family := body.Find(".c-listing-fight__corner-family-name").Text()
	return given + " " + family
}

func trim(str string) string {
	str = strings.TrimSpace(str)
	all := strings.Split(str, "\n")
	for i, _ := range all {
		all[i] = strings.TrimSpace(all[i])
	}
	return strings.Join(all, " ")
}

func NewEvent(doc *goquery.Document) *Event {
	hero := doc.Find(".c-hero__header")
	prefix := trim(hero.Find(".c-hero__headline-prefix h2").Text())
	maineventa := trim(hero.Find(".c-hero__headline > span > span:first-child").Text())
	maineventb := trim(hero.Find(".c-hero__headline > span > span:last-child").Text())
	when := trim(hero.Find(".c-hero__headline-suffix").Text())
	where := trim(doc.Find(".c-hero__text").Text())

	event := Event{
		Title: prefix,
		URL: "",
		Headline: maineventa + " vs " + maineventb,
		Venue: where,
		Start: when,
		Cards: make([]Card, 0, 3),
		Fights: 0,
	}

	if doc.Url != nil {
		event.URL = doc.Url.String()
	}

	// main card, prelims, early prelims
	found := doc.Find("details")
	if found.Length() == 0 {
		event.Cards = append(event.Cards, Card{doc.Selection, "Fight Card", nil})
	} else {
		found.Each(func(n int, sel *goquery.Selection) {
			event.Cards = append(event.Cards, Card{sel, "Unknown Card", nil})
		})
	}

	for idx, _ := range event.Cards {
		summary := event.Cards[idx].root.Find("summary")
		if summary.Length() == 1 {
			event.Cards[idx].Name = summary.Text()
		}

		event.Cards[idx].root.Find("li.l-listing__item").Each(func(n int, sel *goquery.Selection) {
			fight := Fight{
				Blue: Corner{
					Name: getname(sel, "blue"),
					Image: getpict(sel, "blue"),
				},
				Red: Corner{
					Name: getname(sel, "red"),
					Image: getpict(sel, "red"),
				},
			}

			event.Cards[idx].Fights = append(event.Cards[idx].Fights, fight)
			event.Fights += 1
		})
	}

	return &event
}