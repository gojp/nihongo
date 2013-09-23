package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gojp/nihongo/app/routes"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/robfig/revel"
	"log"
	"strings"
)

type App struct {
	*revel.Controller
}

type Gloss struct {
	English  string
	Tags     []string
	Related  []string
	Common   bool
}

type Highlight struct {
	Furigana string
	Japanese string
	Romaji string
	English []string
}

type Word struct {
	Romaji   string
	Common   bool
	Dialects []string
	Fields   []string
	Glosses  []Gloss
	English  []string
	Furigana string
	Japanese string
	Tags     []string
	Pos      []string
}

func search(query string) []Word {
	fmt.Println("Searching for... ", query)
	api.Domain = "localhost"
	searchJson := fmt.Sprintf(`{"query": {"multi_match": {"query": "%s", "fields": ["japanese", "furigana", "romaji", "english"]}}, "highlight": {"fields": {"furigana": {}, "japanese": {}, "romaji": {}, "english": {}}}}`, query)
	out, err := core.SearchRequest(true, "edict", "entry", searchJson, "", 0)
	if err != nil {
		log.Println(err)
	}

	hits := [][]byte{}
	for _, hit := range out.Hits.Hits {
		hits = append(hits, hit.Source)
	}

	wordList := []Word{}
	for _, hit := range hits {
		w := Word{}
		err := json.Unmarshal(hit, &w)
		if err != nil {
			log.Println(err)
		}
		wordList = append(wordList, w)
	}
	return wordList
}

func (a App) Search(query string) revel.Result {
	if len(query) == 0 {
		return a.Redirect(routes.App.Index())
	}
	wordList := search(query)
	return a.Render(wordList)
}

func (c App) Details(query string) revel.Result {
	if len(query) == 0 {
		return c.Redirect(routes.App.Index())
	}
	if (strings.Contains(query, " ")) {
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "-", -1)))
	}
	query = strings.Replace(query, "-", " ", -1)
	wordList := search(query)
	pageTitle := query + " in Japanese"
	return c.Render(wordList, query, pageTitle)
}

func (c App) SearchGet() revel.Result {
	if query, ok := c.Params.Values["q"]; ok && len(query) > 0 {
		return c.Redirect(routes.App.Details(query[0]))
	}
	return c.Redirect(routes.App.Index())
}

func (c App) About() revel.Result {
	return c.Render()
}

func (c App) Index() revel.Result {
	return c.Render()
}
