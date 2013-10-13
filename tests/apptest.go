package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/robfig/revel"
)

type Word struct {
	*models.Word
}

type PopularSearch struct {
	Term string
}

type ScoreWord struct {
	SearchTerm       string
	English          string
	ExpectedPosition int
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

func scoreEnglishPosition(wordList []Word, answer string, expectedPosition int) (score int) {
	score = 0
abc:
	for i, word := range wordList {
		for _, gloss := range word.English {
			if gloss == answer {
				score += 10 - i
				break abc
			}
		}
	}
	return score
}

func (t AppTest) TestSearchResultScores() {
	score := 0

	englishWords := []ScoreWord{
		ScoreWord{"myu-jikku", "music", 0},
		ScoreWord{"test", "test", 0},
	}

	for _, word := range englishWords {
		wordList := getWordList(helpers.Search(word.SearchTerm))
		score += scoreEnglishPosition(wordList, word.English, word.ExpectedPosition)
	}
	finalScore := float64(score*100) / float64(10*len(englishWords))
	fmt.Println("\n\n===================\n Final score is", finalScore, "\n===================\n")

	minimumAllowedScore := 60.0
	t.Assert(finalScore >= minimumAllowedScore)
}

func (t AppTest) After() {
}
