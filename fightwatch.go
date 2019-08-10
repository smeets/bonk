package main
//go:generate templify event.html
//go:generate templify index.html

import (
	"fmt"
	"os"
	"html/template"
	"path"
	"time"
)

func bailout(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func createdatadir(datapath string) {
	_, err := os.Stat(datapath)
	if os.IsNotExist(err) || os.IsExist(err) {
	    err = os.Mkdir(datapath, os.ModeDir)
	    if !os.IsExist(err) {
	    	bailout("mkdata:", err)
	    }
	} else {
		bailout("datadir:", err)
	}

	err = os.Mkdir(path.Join(datapath, "event"), os.ModeDir)
	if !os.IsExist(err) {
		bailout("mkdata:", err)
	}
}

type Page struct {
	*Event
	URL string
}

func main() {
	outdir := "./data"
	createdatadir(outdir)

	eventlinks := findevents()
	pages := make([]Page, 0)

	if len(eventlinks) == 0 {
		fmt.Println("no eventlinks found")
		os.Exit(1)
	}

	evttpl, err := template.New("event").Parse(eventTemplate())
	bailout("event template:", err)

	idxtpl, err := template.New("index").Parse(indexTemplate())
	bailout("index template:", err)

	// Deduplicate some links that seem to come next to each other
	for i := 0; i < len(eventlinks)-1; i++ {
		cur := eventlinks[i]
		nxt := eventlinks[i+1]
		if cur == nxt {
			copy(eventlinks[i:], eventlinks[i+1:])
			eventlinks = eventlinks[:len(eventlinks)-1]
			i -= 1
		}
	}

	for _, eventlink := range eventlinks {
		doc, err := getevent(eventlink)
		if err != nil {
			fmt.Println(eventlink, err)
			continue
		}

		event := NewEvent(doc)
		event.URL = "https://www.ufc.com" + eventlink
		if eventlink[len(eventlink)-1] == '/' {
			eventlink = eventlink[:len(eventlink)-1]
		}

		evtpath := path.Join(outdir, eventlink + ".html")
		evtfile, err := os.Create(evtpath)
		if err != nil {
			fmt.Println("eventgen:", err)
			continue
		}
		defer evtfile.Close()

		evttpl.Execute(evtfile, event)
		pages = append(pages, Page{event,eventlink + ".html"})
		
		fmt.Println("generated", eventlink, evtpath)
	}

	idxfile, err := os.Create(path.Join(outdir, "index.html"))
	if err != nil {
		fmt.Println("eventgen:", idxfile)
		os.Exit(1)
	}
	defer idxfile.Close()

	idxctx := struct {
		Pages []Page
		Created string
	}{
		Pages: pages,
		Created: time.Now().Format("2006-01-02 15:04:05 (-0700 MST)"),
	}
	idxtpl.Execute(idxfile, idxctx)
}