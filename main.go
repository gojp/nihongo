package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/gojp/nihongo/lib/dictionary"
	"github.com/golang/gddo/httputil/header"
)

const (
	title       = "Nihongo.io"
	description = "The world's best Japanese dictionary."
)

// Entry is a dictionary entry
type Entry struct {
	Word       string `json:"word"`
	Furigana   string `json:"furigana"`
	Definition string `json:"definition"`
	Common     bool   `json:"common,omitempty"`
}

var dict dictionary.Dictionary

var tmpl = make(map[string]*template.Template)

func initialize() {
	compileTemplates()

	file, err := os.Open("data/edict2.json.gz")
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

func compileTemplates() {
	t := func(s string) string {
		return "templates/" + s
	}

	tmpl["home.html"] = template.Must(template.ParseFiles(t("home.html"), t("base.html")))
	tmpl["about.html"] = template.Must(template.ParseFiles(t("about.html"), t("base.html")))
}

type templateData struct {
	Search  string  `json:"search"`
	Entries []Entry `json:"entries"`
}

func home(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "/")

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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := map[string]interface{}{
		"json":        string(jsonData),
		"data":        data,
		"title":       title,
		"description": description,
	}

	tmpl["home.html"].ExecuteTemplate(w, "base", m)
}

func search(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "/search")

	// check GET and POST parameters for "text" field
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	text := r.Form.Get("text")

	// if no "text" field is present, we check the URL
	if text == "" {
		text = strings.TrimPrefix(r.URL.Path, "/search/")
	}

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
		log.Println(err)
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
		w.Write(jsonData)
	} else {
		pageTitle := text + " in Japanese | " + title
		description := fmt.Sprintf("Japanese to English for %s", text)
		if len(data.Entries) > 0 {
			e := data.Entries[0]
			description = fmt.Sprintf("%s (%s) - %s", e.Word, e.Furigana, e.Definition)
		}

		m := map[string]interface{}{
			"json":        string(jsonData),
			"data":        data,
			"title":       pageTitle,
			"description": description,
		}
		tmpl["home.html"].ExecuteTemplate(w, "base", m)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	tmpl["about.html"].ExecuteTemplate(w, "base", nil)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var (
		addr string
		dev  bool
	)
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "address to run on")
	flag.BoolVar(&dev, "dev", false, "whether to run with a reduced dictionary (for faster boot times)")
	flag.Parse()

	initialize()

	http.HandleFunc("/", home)
	http.HandleFunc("/search", search)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/about", about)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	log.Printf("Running server on addr %s", addr)
	if dev {
		log.Println("Running in development mode, templates will automatically reload")
	}
	http.ListenAndServe(addr, nil)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
