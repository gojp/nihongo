package models

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
