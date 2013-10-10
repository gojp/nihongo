package controllers

import (
	"encoding/json"
	"github.com/gojp/kana"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/jgraham909/revmgo"
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

type Word struct {
	*models.Word
}

// wrap a string in strong tags
func makeStrong(query string) string {
	return "<strong>" + query + "</strong>"
}

// convert the query to hiragana and katakana. if it's already in
// hiragana or katakana, it will just be the same.
func convertQueryToKana(query string) (hiragana, katakana string) {
	kana := kana.NewKana()
	h := kana.RomajiToHiragana(query)
	k := kana.RomajiToKatakana(query)
	return h, k
}

// Wrap the query in <strong> tags so that we can highlight it in the results
func highlightQuery(w Word, query string) {
	// make regular expression that matches the original query
	re := regexp.MustCompile(`\b` + query + `\b`)
	// convert original query to kana
	h, k := convertQueryToKana(query)
	// make regular expressions that match the hiragana and katakana
	hiraganaRe := regexp.MustCompile(h)
	katakanaRe := regexp.MustCompile(k)
	// wrap the query in strong tags
	queryHighlighted := makeStrong(query)
	hiraganaHighlighted := makeStrong(h)
	katakanaHighlighted := makeStrong(k)

	// if the original input is Japanese, then the original input converted
	// to hiragana and katakana will be equal, so just choose one
	// to highlight so that we only end up with one pair of strong tags
	if hiraganaHighlighted == katakanaHighlighted {
		w.Japanese = hiraganaRe.ReplaceAllString(w.Japanese, hiraganaHighlighted)
	} else {
		// The original input is romaji, so we convert it to hiragana and katakana
		// and highlight both.
		w.Japanese = hiraganaRe.ReplaceAllString(w.Japanese, hiraganaHighlighted)
		w.Japanese = katakanaRe.ReplaceAllString(w.Japanese, katakanaHighlighted)
	}

	// highlight the furigana too, same as above
	w.Furigana = hiraganaRe.ReplaceAllString(w.Furigana, hiraganaHighlighted)
	w.Furigana = katakanaRe.ReplaceAllString(w.Furigana, katakanaHighlighted)
	// highlight the query inside the list of English definitions
	for i, e := range w.English {
		e = re.ReplaceAllString(e, queryHighlighted)
		w.English[i] = e
	}
}

func getWordList(hits [][]byte, query string) (wordList []Word) {
	// highlight queries and build Word object
	for _, hit := range hits {
		w := Word{}
		err := json.Unmarshal(hit, &w)
		if err != nil {
			log.Println(err)
		}
		w.MainEntry = w.Japanese
		highlightQuery(w, query)
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
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "-", -1)))
	}
	// Copy the query so that we maintain the dashes
	// when inserting into MongoDB
	mongoTerm := query
	query = strings.Replace(query, "-", " ", -1)
	hits := helpers.Search(query)
	wordList := getWordList(hits, query)
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
