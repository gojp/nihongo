package edict2

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"
)

type Gloss struct {
	Common  bool
	English string
	Field   *string
	Related []string
	Tags    []string
}

type Entry struct {
	Common    bool
	Dialects  []string
	EntSeq    string `json:"ent_seq"`
	Fields    []string
	Furigana  string
	Glosses   []Gloss
	HasAudio  bool
	Japanese  string
	KanaTags  []string `json:"kana_tags"`
	KanjiTags []string `json:"kanji_tags"`
	Pos       []string
	Tags      []string
}

func Parse(path string) (entries []Entry, err error) {
	file, err := os.Open(path)
	if err != nil {
		return entries, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		return entries, err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var e Entry
		json.Unmarshal(scanner.Bytes(), &e)
		entries = append(entries, e)
	}

	if err := scanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}
