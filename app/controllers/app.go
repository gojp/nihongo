package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/robfig/revel"
	"log"
)

type App struct {
	*revel.Controller
}

func (a App) Search(query string) revel.Result {
	fmt.Println("Searching for... ", query)
	api.Domain = "localhost"
	searchJson := fmt.Sprintf(`{"query": { "multi_match" : {"query" : "%s", "fields" : ["romaji", "furigana", "japanese", "glosses"]}}}`, query)
	out, err := core.SearchRequest(true, "edict", "entry", searchJson, "", 0)
	if err != nil {
		log.Println(err)
	}

	type Word struct {
		Romaji   string
		Common   bool
		Dialects []string
		Fields   []string
		Glosses  []string
		Furigana string
		Japanese string
		Tags     []string
		Pos      []string
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

	return a.Render(wordList)
}

func (c App) Index() revel.Result {
	return c.Render()
}
