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
		Key:        []string{"reading"},
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
	query := bson.M{} // todo: {"reading": search}
	q := collection.Find(query)
	iter := q.Limit(10).Iter()
	iter.All(&wordList)

	fmt.Println(wordList)

	return a.Render(wordList)
}

func (c App) Index() revel.Result {
	return c.Render()
}
