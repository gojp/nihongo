package models

type Gloss struct {
	English      string
	EnglishSplit []string
	Tags         []string
	Related      []string
	Common       bool
}
