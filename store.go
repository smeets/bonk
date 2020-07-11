package main

import (
    "encoding/json"
    "io"
    "os"
)


type Page struct {
    *Event
    URL string
}

type Store struct {
    Events []Page `json:"events"`
}

func LoadStore(path string) *Store {
    file, err := os.Create(path)
    bailout("store:", err)
    defer file.Close()

    store := new(Store)
    if err := json.NewDecoder(file).Decode(store); err != nil && err != io.EOF {
        bailout("store:", err)
    }

    return store
}

func (s *Store) Merge(pages []Page) {
    for _, p := range pages {
        found := false
        for _, x := range s.Events {
            if p.URL == x.URL {
                found = true
                break
            }
        }

        if found {
            continue
        }
        s.Events = append(s.Events, p)
    }
}

func (s *Store) Save(path string) {
    file, err := os.Create(path)
    bailout("save:", err)
    defer file.Close()

    bailout("save:", json.NewEncoder(file).Encode(s))
}
