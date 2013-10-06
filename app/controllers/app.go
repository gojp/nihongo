package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/jgraham909/revmgo"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"regexp"
	"strings"
)

type App struct {
	*revel.Controller
	revmgo.MongoController
}

type Gloss struct {
	English string
	Tags    []string
	Related []string
	Common  bool
}

type Highlight struct {
	Furigana string
	Japanese string
	Romaji   string
	English  []string
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

func (w *Word) highlightQuery(query string) {
	re := regexp.MustCompile(query)
	queryHighlighted := "<strong>" + query + "</strong>"
	w.Japanese = re.ReplaceAllString(w.Japanese, queryHighlighted)
	for i, e := range w.English {
		e = re.ReplaceAllString(e, queryHighlighted)
		w.English[i] = e
	}
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
		w.highlightQuery(query)
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
	if strings.Contains(query, " ") {
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "-", -1)))
	}
	query = strings.Replace(query, "-", " ", -1)
	wordList := search(query)
	pageTitle := query + " in Japanese"

	// log this call in mongo
	collection := c.MongoSession.DB("greenbook").C("hits")
	_, err := collection.Upsert(bson.M{"term": query}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		// mongo failed to log, but who cares
	}

	index := mgo.Index{
		Key:        []string{"count"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	collection.EnsureIndex(index)

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
	collection := c.MongoSession.DB("greenbook").C("hits")
	q := collection.Find(nil).Sort("-count")

	termList := []models.SearchTerm{}
	iter := q.Limit(10).Iter()
	iter.All(&termList)

	return c.Render(termList)
}
