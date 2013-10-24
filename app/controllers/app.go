package controllers

import (
	"encoding/json"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/robfig/revel"
	"log"
	"strings"
)

type App struct {
	*revel.Controller
}

type Word struct {
	*models.Word
}

type PopularSearch struct {
	Term string
}

func getWordList(hits [][]byte, query string) (wordList []Word) {
	// highlight queries and build Word object
	for _, hit := range hits {
		w := Word{}
		err := json.Unmarshal(hit, &w)
		if err != nil {
			log.Println(err)
		}
		w.HighlightQuery(query)
		wordList = append(wordList, w)
	}
	return wordList
}

func (a App) Search(query string) revel.Result {
	if len(query) == 0 {
		return a.Redirect(routes.App.Index())
	}
	hits := helpers.Search(query)
	wordList := getWordList(hits, query)
	return a.Render(wordList)
}

func (c App) Details(query string) revel.Result {
	if len(query) == 0 {
		return c.Redirect(routes.App.Index())
	}
	if strings.Contains(query, " ") {
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "_", -1)))
	}

	query = strings.Replace(query, "_", " ", -1)
	hits := helpers.Search(query)
	wordList := getWordList(hits, query)
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

	// get the popular searches
	// collection := c.MongoSession.DB("greenbook").C("hits")
	// q := collection.Find(nil).Sort("-count")

	// termList := []models.SearchTerm{}
	// iter := q.Limit(10).Iter()
	// iter.All(&termList)

	termList := []PopularSearch{
		PopularSearch{"今日は"},
		PopularSearch{"kanji"},
		PopularSearch{"amazing"},
		PopularSearch{"かんじ"},
		PopularSearch{"莞爾"},
		PopularSearch{"天真流露"},
		PopularSearch{"funny"},
		PopularSearch{"にほんご"},
	}

	return c.Render(termList)
}
