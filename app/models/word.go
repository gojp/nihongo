package models

import (
	"github.com/gojp/nihongo/app/helpers"
	"regexp"
	"strings"
)

type Word struct {
	Romaji     string
	Common     bool
	Dialects   []string
	Fields     []string
	Glosses    []Gloss
	English    []string
	EnglishHL  []string // highlighted english
	Furigana   string
	FuriganaHL string // highlighted furigana
	Japanese   string
	JapaneseHL string // highlighted japanese
	Tags       []string
	Pos        []string
}

// Wrap the query in <strong> tags so that we can highlight it in the results
func (w *Word) HighlightQuery(query string) {
	// make regular expression that matches the original query
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(query) + `\b`)
	// convert original query to kana
	h, k := helpers.ConvertQueryToKana(query)
	// wrap the query in strong tags
	queryHighlighted := helpers.MakeStrong(query)
	hiraganaHighlighted := helpers.MakeStrong(h)
	katakanaHighlighted := helpers.MakeStrong(k)

	// if the original input is Japanese, then the original input converted
	// to hiragana and katakana will be equal, so just choose one
	// to highlight so that we only end up with one pair of strong tags
	w.JapaneseHL = strings.Replace(w.Japanese, h, hiraganaHighlighted, -1)
	if h != k {
		w.JapaneseHL = strings.Replace(w.JapaneseHL, k, katakanaHighlighted, -1)
	}

	// highlight the furigana too, same as above
	w.FuriganaHL = strings.Replace(w.Furigana, h, hiraganaHighlighted, -1)
	if h != k {
		w.FuriganaHL = strings.Replace(w.FuriganaHL, k, katakanaHighlighted, -1)
	}

	// highlight the query inside the list of English definitions
	w.EnglishHL = []string{}
	for _, e := range w.English {
		e = re.ReplaceAllString(e, queryHighlighted)
		w.EnglishHL = append(w.EnglishHL, e)
	}
}
