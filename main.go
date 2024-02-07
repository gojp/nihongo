package main

import (
	"compress/gzip"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gojp/nihongo/lib/dictionary"
	"github.com/golang/gddo/httputil/header"
)

const (
	title       = "Nihongo.io"
	description = "The world's best Japanese dictionary."
)

var (
	//go:embed templates/*
	content embed.FS

	//go:embed data/*
	data embed.FS

	//go:embed static/*
	static embed.FS
)

// Entry is a dictionary entry
type Entry struct {
	Word       string `json:"word"`
	Furigana   string `json:"furigana"`
	Definition string `json:"definition"`
	Common     bool   `json:"common,omitempty"`
}

var dict dictionary.Dictionary

func initialize() {
	file, err := data.Open("data/edict2.json.gz")
	if err != nil {
		log.Fatal("Could not load edict2.json.gz: ", err)
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal("Could not create reader: ", err)
	}

	dict, err = dictionary.Load(reader)
	if err != nil {
		log.Fatal("Could not load dictionary: ", err)
	}
}

type templateData struct {
	Search  string  `json:"search"`
	Entries []Entry `json:"entries"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// defer timeTrack(time.Now(), "/")

	if r.URL.Path[1:] != "" {
		http.NotFound(w, r)
		return
	}

	data := templateData{
		Entries: []Entry{},
		Search:  "",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := map[string]interface{}{
		"json":        template.HTML(string(jsonData)),
		"data":        data,
		"title":       title,
		"description": description,
	}

	t, err := template.ParseFS(content, "templates/home.html")
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "home.html", m)
	if err != nil {
		log.Println("ERROR:", err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	// check GET and POST parameters for "text" field
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	text := r.Form.Get("text")
	trackText := fmt.Sprintf("/search?text=%s", text)

	// if no "text" field is present, we check the URL
	if text == "" {
		text = strings.TrimPrefix(r.URL.Path, "/search/")
		trackText = fmt.Sprintf("/search/%s", text)
	}

	defer timeTrack(time.Now(), trackText)

	// if we still don't have text, we redirect to the home page
	if text == "" {
		log.Println("Redirecting to home")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// get the entries that match our text
	entries := []Entry{}
	results := dict.Search(text, 10)
	for _, r := range results {
		var defs []string
		for _, g := range r.Glosses {
			defs = append(defs, g.English)
		}
		entries = append(entries, Entry{
			Word:       r.Japanese,
			Furigana:   r.Furigana,
			Definition: strings.Join(defs, "; "),
			Common:     r.Common,
		})
	}

	data := templateData{
		Search:  text,
		Entries: entries,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isXMLHTTP := r.Header.Get("X-Requested-With") == "XMLHttpRequest"
	accepts := header.ParseAccept(r.Header, "Accept")
	wantsJSON, wantsHTML := 0.0, 0.0
	for _, acc := range accepts {
		switch acc.Value {
		case "text/json", "application/json":
			wantsJSON = acc.Q
		case "text/html":
			wantsHTML = acc.Q
		}
	}
	if isXMLHTTP || wantsJSON > wantsHTML {
		w.Header().Set("Content-Type", "application/json")
		b := []byte(template.HTML(string(jsonData)))
		w.Write(b)

		return
	}

	pageTitle := text + " in Japanese | " + title
	description := fmt.Sprintf("Japanese to English for %s", text)
	if len(data.Entries) > 0 {
		e := data.Entries[0]
		description = fmt.Sprintf("%s (%s) - %s", e.Word, e.Furigana, e.Definition)
	}

	m := map[string]interface{}{
		"json":        template.HTML(string(jsonData)),
		"data":        data,
		"title":       pageTitle,
		"description": description,
	}

	t, err := template.ParseFS(content, "templates/home.html")
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "home.html", m)
	if err != nil {
		log.Println("ERROR:", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		log.Println("ERROR:", err)
	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "address to run on")
	flag.Parse()

	initialize()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", search)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/about", aboutHandler)
	http.Handle("/static/", http.FileServer(http.FS(static)))

	log.Printf("Running on %s ...", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
