package models

import (
	"labix.org/v2/mgo/bson"
)

type WordModel struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Info     string
	Glosses  []string
	Furigana string
	Japanese string
	Romaji   string
	Common   bool
}

type SearchTerm struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Term  string
	Count int
}
