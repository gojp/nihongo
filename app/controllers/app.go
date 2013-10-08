package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gojp/kana"
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
	Romaji    string
	Common    bool
	Dialects  []string
	Fields    []string
	Glosses   []Gloss
	English   []string
	Furigana  string
	Japanese  string
	MainEntry string
	Tags      []string
	Pos       []string
}

// wrap a string in strong tags
func makeStrong(query string) string {
	return "<strong>" + query + "</strong>"
}

// convert the query to hiragana and katakana. if it's already in
// hiragana or katakana, it will just be the same.
func convertQueryToKana(query string) (hiragana, katakana string) {
	kana := kana.NewKana()
	h := kana.Romaji_to_hiragana(query)
	k := kana.Romaji_to_katakana(query)
	return h, k
}

// Wrap the query in <strong> tags so that we can highlight it in the results
func (w *Word) highlightQuery(query string) {
	// make regular expression that matches the original query
	re := regexp.MustCompile(`\b` + query + `\b`)
	// convert original query to kana
	h, k := convertQueryToKana(query)
	// make regular expressions that match the hiragana and katakana
	hiraganaRe := regexp.MustCompile(h)
	katakanaRe := regexp.MustCompile(k)
	// wrap the query in strong tags
	queryHighlighted := makeStrong(query)
	katakanaHighlighted := makeStrong(k)
	hiraganaHighlighted := makeStrong(h)
	// if the query is originally in Japanese, highlight it
	w.Japanese = re.ReplaceAllString(w.Japanese, queryHighlighted)
	// highlight the katakana or hiragana that has been converted from romaji
	w.Japanese = hiraganaRe.ReplaceAllString(w.Japanese, hiraganaHighlighted)
	w.Japanese = katakanaRe.ReplaceAllString(w.Japanese, katakanaHighlighted)
	// highlight the furigana too, same as above
	w.Furigana = hiraganaRe.ReplaceAllString(w.Furigana, hiraganaHighlighted)
	w.Furigana = katakanaRe.ReplaceAllString(w.Furigana, katakanaHighlighted)
	// highlight the query inside the list of English definitions
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
		w.MainEntry = w.Japanese
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
    // Copy the query so that we maintain the dashes
    // when inserting into MongoDB
	mongoTerm := query
	query = strings.Replace(query, "-", " ", -1)
	wordList := search(query)
	pageTitle := query + " in Japanese"

	// log this call in mongo
	collection := c.MongoSession.DB("greenbook").C("hits")
	_, err := collection.Upsert(bson.M{"term": mongoTerm}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		log.Println("DEBUG: mongo failed to upsert count: " + err.Error())
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
