package controllers

import (
	"fmt"
	//"github.com/gojp/nihongo/app/models"
	//"github.com/mattbaird/elastigo/api"
	//"github.com/mattbaird/elastigo/search"
	//"github.com/mattbaird/elastigo/core"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (a App) Search(query string) revel.Result {
	//api.Domain = "localhost"
	//out, err := search.Search("edict").Type("entry").Size("100").Search(query).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(out)

	fmt.Println("Searching for... ", query)

	// Index - not the best place for this, but okay for now...
	//index := mgo.Index{
	//	Key:        []string{"romaji", "furigana", "japanese", "glosses"},
	//	Unique:     false,
	//	DropDups:   false,
	//	Background: true,
	//	Sparse:     true,
	//}

	//err := collection.EnsureIndex(index)

	//if err != nil {
	//	panic("Database connection failed")
	//}

	//wordList := []models.Word{}
	//query := bson.M{"$or": []bson.M{
	//	bson.M{"romaji": bson.RegEx{".*" + search + ".*", "i"}},
	//	bson.M{"furigana": search},
	//	bson.M{"japanese": search},
	//}}
	//q := collection.Find(query).Sort("-common", "furigana")
	//iter := q.Limit(100).Iter()
	//iter.All(&wordList)

	return a.Render()
}

func (c App) Index() revel.Result {
	return c.Render()
}
