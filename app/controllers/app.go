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
	h := kana.RomajiToHiragana(query)
	k := kana.RomajiToKatakana(query)
	return h, k
}

// Wrap the query in <strong> tags so that we can highlight it in the results
func (w *Word) highlightQuery(query string) {
	// make regular expression that matches the original query
	re := regexp.MustCompile(`\b` + query + `\b`)
	// convert original query to kana
	h, k := convertQueryToKana(query)
	// wrap the query in strong tags
	queryHighlighted := makeStrong(query)
	hiraganaHighlighted := makeStrong(h)
	katakanaHighlighted := makeStrong(k)

	// if the original input is Japanese, then the original input converted
	// to hiragana and katakana will be equal, so just choose one
	// to highlight so that we only end up with one pair of strong tags
	if hiraganaHighlighted == katakanaHighlighted {
		w.Japanese = strings.Replace(w.Japanese, h, hiraganaHighlighted, -1)
	} else {
		// The original input is romaji, so we convert it to hiragana and katakana
		// and highlight both.
		w.Japanese = strings.Replace(w.Japanese, h, hiraganaHighlighted, -1)
		w.Japanese = strings.Replace(w.Japanese, k, katakanaHighlighted, -1)
	}

	// highlight the furigana too, same as above
	w.Furigana = strings.Replace(w.Furigana, h, hiraganaHighlighted, -1)
	w.Furigana = strings.Replace(w.Furigana, k, katakanaHighlighted, -1)
	// highlight the query inside the list of English definitions
	for i, e := range w.English {
		e = re.ReplaceAllString(e, queryHighlighted)
		w.English[i] = e
	}
}

func search(query string) []Word {
	fmt.Println("Searching for... ", query)
	api.Domain = "localhost"

	kana := kana.NewKana()

	isLatin := kana.IsLatin(query)
	isKana := kana.IsKana(query)
	isKanji := kana.IsKanji(query)
	fmt.Println(isLatin, isKana, isKanji)

	// convert to hiragana and katakana
	romaji := kana.KanaToRomaji(query)

	// handle different types of input differently:
	matches := []string{}
	if isKana {
		// add boost for exact-matching kana
		matches = append(matches, fmt.Sprintf(`
		{"match" : 
			{
				"furigana" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 5.0
				}
			}
		}`, query))

		// also look for romaji version in case
		matches = append(matches, fmt.Sprintf(`
		{"match" : 
			{
				"romaji" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 2.0
				}
			}
		}`, romaji))
	}
	if !isLatin {
		matches = append(matches, fmt.Sprintf(`
		{"match" : 
			{
				"japanese" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 10.0
				}
			}
		}`, query))
	} else {
		// add romaji search term
		matches = append(matches, fmt.Sprintf(`
		{"match" : 
			{
				"romaji" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 3.0
				}
			}
		}`, query))

		// add english search term
		matches = append(matches, fmt.Sprintf(`
		{"match" : 
			{
				"english" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 5.0
				}
			}
		}`, query))
	}

	searchJson := fmt.Sprintf(`
		{"query": 
			{"bool": 
				{
				"should":
					[` + strings.Join(matches, ",") + `],
				"minimum_should_match" : 0,
				"boost": 2.0
				}
			}
		}`)

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
