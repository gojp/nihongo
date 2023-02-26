package dictionary

import (
	"container/heap"
	"io"
	"sort"
	"strings"

	"github.com/gojp/kana"
	"github.com/gojp/nihongo/edict2"
)

type EntryID uint64

type Entry struct {
	edict2.Entry
	ID EntryID
}

type Dictionary struct {
	entries  map[EntryID]Entry
	japanese *RadixTree
	furigana *RadixTree
	english  *InvertedIndex
}

func newEntry(entry *edict2.Entry, id uint64) *Entry {
	return &Entry{
		Entry: *entry,
		ID:    EntryID(id),
	}
}

func cleanWord(w string) string {
	return strings.ToLower(strings.Trim(w, ",.()|`'\"!"))
}

func Load(r io.Reader) (Dictionary, error) {
	d := Dictionary{}
	d.entries = map[EntryID]Entry{}
	d.japanese = NewRadixTree()
	d.furigana = NewRadixTree()
	d.english = NewInvertedIndex(30)

	edict := edict2.New(r)
	var i uint64
	for edict.Scan() {
		i++
		err := edict.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			return d, err
		}
		e := newEntry(edict.Entry(), i)
		d.entries[e.ID] = *e
		d.japanese.Insert(e.Japanese, e.ID)
		d.furigana.Insert(e.Furigana, e.ID)

		for _, gloss := range e.Glosses {
			words := strings.Split(gloss.English, " ")
			for _, w := range words {
				d.english.Insert(cleanWord(w), e.ID, 1.0/float64(len(words)))
			}
		}

	}
	if err := edict.Err(); err != nil {
		return d, err
	}

	return d, nil
}

// Get fetches the entry with the given ID, and returns it.
func (d Dictionary) Get(id EntryID) (e Entry, found bool) {
	e, found = d.entries[id]
	return
}

type ByCommon []Entry

func (a ByCommon) Len() int           { return len(a) }
func (a ByCommon) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommon) Less(i, j int) bool { return a[i].Common && !a[j].Common }

// Search takes a search string provided by the user, and returns a matching
// Entry slice with at most `limit` number of entries.
func (d Dictionary) Search(s string, limit int) (results []Entry) {
	results = []Entry{}
	resultsMap := map[EntryID]bool{}

	word := cleanWord(s)

	appendResults := func(f func(word string, max int) []EntryID, word string, max int) {
		if entryIDs := f(word, max); entryIDs != nil {
			for _, eid := range entryIDs {
				// some entries have the same Japanese and Furigana fields, so we should
				// only add those to the results slice once
				if _, found := resultsMap[eid]; found {
					continue
				}

				resultsMap[eid] = true

				results = append(results, d.entries[eid])
			}
		}
	}

	appendResults(d.japanese.FindWordsWithPrefix, word, 5)
	appendResults(d.furigana.FindWordsWithPrefix, word, 5)

	if kana.IsLatin(word) {
		hira := kana.RomajiToHiragana(word)
		kata := kana.RomajiToKatakana(word)

		appendResults(d.furigana.FindWordsWithPrefix, hira, 5)
		appendResults(d.furigana.FindWordsWithPrefix, kata, 5)
	}

	// build a priority queue of relevant entries for english search terms,
	// using our inverted index, and pull out the top 10
	words := strings.Split(s, " ")
	scores := map[EntryID]float64{}
	for i, w := range words {
		// limit to 10 words, to keep the request time bounded
		if i > 10 {
			break
		}

		cw := cleanWord(w)
		r := d.english.Get(cw)
		if r == nil {
			continue
		}

		for i := range r {
			scores[r[i].id] += r[i].score
			for _, w2 := range words {
				cw2 := cleanWord(w2)
				if w != w2 {
					if d.english.Test(cw2, r[i].id) {
						scores[r[i].id] *= 10 // for now double score without verifying bloom-filter correctness
					} else {
						scores[r[i].id] /= 10
					}
				}
			}
		}
	}

	pq := make(PriorityQueue, len(scores))
	i := 0
	for id, priority := range scores {
		pq[i] = &Item{
			id:       id,
			priority: priority,
			index:    i,
		}
		i++
	}

	heap.Init(&pq)
	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 && len(results) < 10 {
		item := heap.Pop(&pq).(*Item)
		results = append(results, d.entries[item.id])
	}

	sort.Sort(ByCommon(results))
	return
}
