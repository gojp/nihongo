package tests

import (
	"github.com/gojp/nihongo/app/models"
	"github.com/robfig/revel"
)

type ModelsTest struct {
	revel.TestSuite
}

func (s *ModelsTest) Before() {
}

func (s *ModelsTest) After() {
}

func (s *ModelsTest) TestHighlightQuery() {
	// some basic checks
	w := models.Word{
		English:  []string{"test"},
		Furigana: "テスト",
		Japanese: "テスト",
	}
	w.HighlightQuery("tesuto")
	s.Assert(w.English[0] == "test")
	s.Assert(w.EnglishHL[0] == "test")
	s.Assert(w.Furigana == "テスト")
	s.Assert(w.FuriganaHL == "<strong>テスト</strong>")
	s.Assert(w.Japanese == "テスト")
	s.Assert(w.JapaneseHL == "<strong>テスト</strong>")

	w.HighlightQuery("テスト")
	s.Assert(w.English[0] == "test")
	s.Assert(w.EnglishHL[0] == "test")
	s.Assert(w.Furigana == "テスト")
	s.Assert(w.FuriganaHL == "<strong>テスト</strong>")
	s.Assert(w.Japanese == "テスト")
	s.Assert(w.JapaneseHL == "<strong>テスト</strong>")

	w.HighlightQuery("test")
	s.Assert(w.English[0] == "test")
	s.Assert(w.EnglishHL[0] == "<strong>test</strong>")
	s.Assert(w.Furigana == "テスト")
	s.Assert(w.FuriganaHL == "テスト")
	s.Assert(w.Japanese == "テスト")
	s.Assert(w.JapaneseHL == "テスト")
}
