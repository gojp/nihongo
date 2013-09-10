package models

import (
	"labix.org/v2/mgo/bson"
)

type Word struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Info     string
	Glosses  []string
	Furigana string
	Japanese string
	Romaji   string
	Common   bool
}
