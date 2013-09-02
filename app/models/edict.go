package models

import (
	"labix.org/v2/mgo/bson"
)

type Word struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Info    string
	Reading string
}
