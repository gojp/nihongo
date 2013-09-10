package controllers

import (
	"fmt"
	"github.com/gojp/greenbook/app/models"
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type App struct {
	*revel.Controller
	revmgo.MongoController
}

func (a App) Search(search string) revel.Result {
	collection := a.MongoSession.DB("greenbook").C("edict")

	fmt.Println("Searching for... ", search)

	// Index - not the best place for this, but okay for now...
	index := mgo.Index{
		Key:        []string{"romaji", "furigana", "japanese", "glosses"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)

	if err != nil {
		panic("Database connection failed")
	}

	wordList := []models.Word{}
	query := bson.M{"$or": []bson.M{
		bson.M{"romaji": bson.RegEx{".*" + search + ".*", "i"}},
		bson.M{"furigana": search},
		bson.M{"japanese": search},
	}}
	q := collection.Find(query).Sort("-common", "furigana")
	iter := q.Limit(100).Iter()
	iter.All(&wordList)

	return a.Render(wordList)
}

func (c App) Index() revel.Result {
	return c.Render()
}
