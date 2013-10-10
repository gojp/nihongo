package models

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
