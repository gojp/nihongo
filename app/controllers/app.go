package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"strings"
)

type App struct {
	*revel.Controller
	revmgo.MongoController
}

type Word struct {
	*models.Word
}

type PopularSearch struct {
	Term string
}

func getWordList(hits [][]byte, query string) (wordList []Word) {
	// highlight queries and build Word object
	for _, hit := range hits {
		w := Word{}
		err := json.Unmarshal(hit, &w)
		if err != nil {
			log.Println(err)
		}
		w.HighlightQuery(query)
		wordList = append(wordList, w)
	}
	return wordList
}

func (a App) Search(query string) revel.Result {
	if len(query) == 0 {
		return a.Redirect(routes.App.Index())
	}
	hits := helpers.Search(query)
	wordList := getWordList(hits, query)
	return a.Render(wordList)
}

func (c App) Details(query string) revel.Result {
	if len(query) == 0 {
		return c.Redirect(routes.App.Index())
	}
	if strings.Contains(query, " ") {
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "_", -1)))
	}

	query = strings.Replace(query, "_", " ", -1)
	hits := helpers.Search(query)
	wordList := getWordList(hits, query)
	pageTitle := query + " in Japanese"

	return c.Render(wordList, query, pageTitle)
}

func (c App) SearchGet() revel.Result {
	if query, ok := c.Params.Values["q"]; ok && len(query) > 0 {
		return c.Redirect(routes.App.Details(query[0]))
	}
	return c.Redirect(routes.App.Index())
}

func (c App) About() revel.Result {
	return c.Render()
}

func addUser(collection *mgo.Collection, email, password string) {
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		log.Panic(err)
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)

	err = collection.Insert(&models.User{Email: email, Password: string(bcryptPassword)})

	if err != nil {
		log.Panic(err)
	}
}

func (c App) Register() revel.Result {
	title := "Register"
	return c.Render(title)
}

func (c App) SaveUser(user models.User) revel.Result {
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}

	collection := c.MongoSession.DB("greenbook").C("users")
	addUser(collection, user.Email, user.Password)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.Email)
	return c.Redirect(routes.App.Index())
}

func (c App) getUser(email string) *models.User {
	users := c.MongoSession.DB("greenbook").C("users")
	result := models.User{}
	users.Find(bson.M{"email": email}).One(&result)
	return &result
}

func (c App) Login(email, password string) revel.Result {
	user := c.getUser(email)
	if user != nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {
			c.Session["user"] = email
			c.Flash.Success("Welcome, " + email)
			return c.Redirect(routes.App.Index())
		}
	}

	c.Flash.Out["email"] = email
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func (c App) Index() revel.Result {

	// get the popular searches
	// collection := c.MongoSession.DB("greenbook").C("hits")
	// q := collection.Find(nil).Sort("-count")

	// termList := []models.SearchTerm{}
	// iter := q.Limit(10).Iter()
	// iter.All(&termList)

	termList := []PopularSearch{
		PopularSearch{"今日は"},
		PopularSearch{"kanji"},
		PopularSearch{"amazing"},
		PopularSearch{"かんじ"},
		PopularSearch{"莞爾"},
		PopularSearch{"天真流露"},
		PopularSearch{"funny"},
		PopularSearch{"にほんご"},
	}

	return c.Render(termList)
}
