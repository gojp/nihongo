package tests

import (
	"encoding/json"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/robfig/revel"
)

type Word struct {
	*models.Word
}

func getWordList(hits [][]byte) (wordList []Word) {
	// highlight queries and build Word object
	for _, hit := range hits {
		w := Word{}
		json.Unmarshal(hit, &w)
		wordList = append(wordList, w)
	}
	return wordList
}

type AppTest struct {
	revel.TestSuite
}

func (t AppTest) Before() {
	println("Set up")
}

func (t AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatHelloSearchPageWorks() {
	t.Get("/hello")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatKonnichiwaSearchPageWorks() {
	t.Get("/今日は")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatDoublePlusSearchWorks() {
	t.Get("/今日は++")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestSearchResults() {
	// some basic checks
	wordList := getWordList(helpers.Search("hello"))
	t.Assert(wordList[0].Japanese == "今日は")

	wordList = getWordList(helpers.Search("kokoro"))
	t.Assert(wordList[0].Japanese == "心")

	wordList = getWordList(helpers.Search("心"))
	t.Assert(wordList[0].Japanese == "心")
}

func (t AppTest) TestSearchResultScores() {
	wordList := getWordList(helpers.Search("myu-jikku"))
	t.Assert(wordList[0].English[0] == "music")
}

func (t AppTest) After() {
	println("Tear down")
}
