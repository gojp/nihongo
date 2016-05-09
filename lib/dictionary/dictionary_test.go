package dictionary

import (
	"compress/gzip"
	"fmt"
	"os"
	"testing"
)

func loadDict() (*Dictionary, error) {
	file, err := os.Open("../../data/edict2.json.gz")
	if err != nil {
		return nil, fmt.Errorf("could not load edict2.json.gz file: %s", err)
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("could not create reader: %s", err)
	}
	dict, err := Load(reader)
	if err != nil {
		return nil, err
	}

	return &dict, nil
}

var exactSearches = []struct {
	word string
	want int
}{
	{"むかつく", 1},
	{"龜兒子龜兒子", 0},
}

func TestExactMatch(t *testing.T) {
	d, err := loadDict()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range exactSearches {
		results := d.Search(s.word, 10)
		if len(results) != s.want {
			t.Errorf("d.Search(%q) returned %d results, want %d", s.word, len(results), s.want)
		}
	}
}
