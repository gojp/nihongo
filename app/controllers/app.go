package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gojp/nihongo/app/helpers"
	"github.com/gojp/nihongo/app/models"
	"github.com/gojp/nihongo/app/routes"
	"github.com/jgraham909/revmgo"
	"github.com/revel/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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

func (c App) connected() *models.User {
	if c.RenderArgs["email"] != nil {
		return c.RenderArgs["email"].(*models.User)
	}
	if email, ok := c.Session["email"]; ok {
		return c.getUser(email)
	}
	return nil
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
		return a.Redirect(App.Index)
	}
	hits := helpers.Search(query)
	fuzzy := false

	if len(hits) == 0 {
		// no hits, so we make suggestions ("did you mean...")
		hits = helpers.FuzzySearch(query)
		fuzzy = true
	}
	wordList := getWordList(hits, query)

	foundSomething := len(wordList)

	return a.Render(wordList, fuzzy, foundSomething, query)
}

func (c App) Details(query string) revel.Result {
	if len(query) == 0 {
		return c.Redirect(App.Index)
	}
	if strings.Contains(query, " ") {
		return c.Redirect(routes.App.Details(strings.Replace(query, " ", "_", -1)))
	}

	query = strings.Replace(query, "_", " ", -1)
	hits := helpers.Search(query)
	fuzzy := false

	if len(hits) == 0 {
		// no hits, so we make suggestions ("did you mean...")
		hits = helpers.FuzzySearch(query)
		fuzzy = true
	}

	wordList := getWordList(hits, query)
	pageTitle := query + " in Japanese"

	description := ""
	if len(wordList) > 0 {
		w := wordList[0].Word
		if w.Japanese != w.Furigana {
			description += fmt.Sprintf("%s [%s] - ", w.Japanese, w.Furigana)
		} else {
			description += fmt.Sprintf("%s - ", w.Japanese)
		}
		description += strings.Join(w.English, ", ")
	}

	foundSomething := len(wordList)
	user := c.connected()
	return c.Render(wordList, query, pageTitle, user, description, fuzzy, foundSomething)
}

func (c App) SearchGet() revel.Result {
	if query, ok := c.Params.Values["q"]; ok && len(query) > 0 {
		return c.Redirect(routes.App.Details(query[0]))
	}
	return c.Redirect(App.Index)
}

func (c App) About() revel.Result {
	return c.Render()
}

func (c App) Resources() revel.Result {
	return c.Render()
}

func (a App) SavePhrase(phrase string) revel.Result {
	if len(phrase) == 0 || a.connected() == nil {
		return a.Redirect(App.Index)
	}
	user := a.connected()
	user.Words = append(user.Words, phrase)

	// todo: should be in model save function or the like
	collection := a.MongoSession.DB("greenbook").C("users")
	err := collection.Update(bson.M{"email": user.Email}, bson.M{"$set": bson.M{"words": user.Words}})

	if err != nil {
		log.Panic(err)
	}

	return a.RenderJson(bson.M{"result": "ok"})
}

func addUser(collection *mgo.Collection, email, password string) error {
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
		return err
	}

	return nil
}

func (c App) Register() revel.Result {
	title := "Register"
	return c.Render(title)
}

func (c App) LoginPage() revel.Result {
	title := "Login"
	return c.Render(title)
}

func (c App) Profile() revel.Result {
	user := c.connected()
	wordList := user.Words
	return c.Render(wordList)
}

func (c App) SaveUser(user models.User) revel.Result {
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Register)
	}

	collection := c.MongoSession.DB("greenbook").C("users")
	err := addUser(collection, user.Email, user.Password)

	if err != nil {
		if mgo.IsDup(err) {
			c.Flash.Error("We're sorry, but a user with this email address already exists.")
			return c.Redirect(App.Register)
		} else {
			c.Flash.Error("We're sorry, but we are experiencing difficulties adding users to the system. Please try again later.")
			log.Println("ERROR: could not add user: ", err)
			return c.Redirect(App.Register)
		}
	}

	c.Session["email"] = user.Email
	c.Flash.Success("Thanks for joining Nihongo.io. よろしくお願いします!")
	return c.Redirect(App.Index)
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
			c.Session["email"] = email
			c.Flash.Success("Welcome, " + email)
			return c.Redirect(App.Index)
		}
	}

	c.Flash.Out["email"] = email
	c.Flash.Error("Login failed")
	return c.Redirect(App.Index)
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(App.Index)
}

func (c App) Index() revel.Result {

	// get the popular searches
	// collection := c.MongoSession.DB("greenbook").C("hits")
	// q := collection.Find(nil).Sort("-count")

	// termList := []models.SearchTerm{}
	// iter := q.Limit(10).Iter()
	// iter.All(&termList)

	termList := []PopularSearch{
		PopularSearch{"こんにちは"},
		PopularSearch{"kanji"},
		PopularSearch{"amazing"},
		PopularSearch{"かんじ"},
		PopularSearch{"莞爾"},
		PopularSearch{"天真流露"},
		PopularSearch{"funny"},
		PopularSearch{"にほんご"},
	}
	user := c.connected()
	return c.Render(termList, user)
}
