package edict2

import (
	"bufio"
	"encoding/json"
	"io"
)

type EDict struct {
	*bufio.Scanner
	TokenType int
	entry     *Entry
}

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

func New(r io.Reader) *EDict {
	s := bufio.NewScanner(r)
	edict := &EDict{
		Scanner: s,
	}
	return edict
}

func (edict *EDict) NextEntry() error {
	e, err := parseEntry(edict.Bytes())
	if err != nil {
		return err
	}
	edict.entry = e

	return nil
}

func (e *EDict) Entry() *Entry {
	return e.entry
}

func parseEntry(entry []byte) (*Entry, error) {
	var e Entry
	json.Unmarshal(entry, &e)

	return &e, nil
}
