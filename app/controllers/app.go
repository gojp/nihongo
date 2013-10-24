package controllers

import (
	"encoding/json"
	"github.com/gojp/kana"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/robfig/revel"
	"log"
	"regexp"
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

// convert the query to hiragana and katakana. if it's already in
// hiragana or katakana, it will just be the same.
func convertQueryToKana(query string) (hiragana, katakana string) {
	h := kana.RomajiToHiragana(query)
	k := kana.RomajiToKatakana(query)
	return h, k
}

// Wrap the query in <strong> tags so that we can highlight it in the results
func highlightQuery(w Word, query string) {
	// make regular expression that matches the original query
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(query) + `\b`)
	// convert original query to kana
	h, k := convertQueryToKana(query)
	// wrap the query in strong tags
	queryHighlighted := helpers.MakeStrong(query)
	hiraganaHighlighted := helpers.MakeStrong(h)
	katakanaHighlighted := helpers.MakeStrong(k)

	// if the original input is Japanese, then the original input converted
	// to hiragana and katakana will be equal, so just choose one
	// to highlight so that we only end up with one pair of strong tags
	if hiraganaHighlighted == katakanaHighlighted {
		w.JapaneseHL = strings.Replace(w.Japanese, h, hiraganaHighlighted, -1)
	} else {
		// The original input is romaji, so we convert it to hiragana and katakana
		// and highlight both.
		w.JapaneseHL = strings.Replace(w.Japanese, h, hiraganaHighlighted, -1)
		w.JapaneseHL = strings.Replace(w.JapaneseHL, k, katakanaHighlighted, -1)
	}

	// highlight the furigana too, same as above
	w.FuriganaHL = strings.Replace(w.Furigana, h, hiraganaHighlighted, -1)
	w.FuriganaHL = strings.Replace(w.FuriganaHL, k, katakanaHighlighted, -1)
	// highlight the query inside the list of English definitions
	w.EnglishHL = []string{}
	for _, e := range w.English {
		e = re.ReplaceAllString(e, queryHighlighted)
		w.EnglishHL = append(w.EnglishHL, e)
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
