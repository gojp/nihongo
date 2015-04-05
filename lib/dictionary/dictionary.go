package dictionary

import (
	"container/heap"
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

func Load(path string) (Dictionary, error) {
	d := Dictionary{}
	d.entries = map[EntryID]Entry{}
	d.japanese = NewRadixTree()
	d.furigana = NewRadixTree()
	d.english = NewInvertedIndex(30)

	entries, err := edict2.Parse(path)
	if err != nil {
		return d, err
	}
	for i, entry := range entries {
		e := newEntry(&entry, uint64(i))
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

	if entryIDs := d.japanese.Get(cleanWord(s)); entryIDs != nil {
		for _, eid := range entryIDs {
			results = append(results, d.entries[eid])
		}
	}

	if entryIDs := d.furigana.Get(cleanWord(s)); entryIDs != nil {
		for _, eid := range entryIDs {
			results = append(results, d.entries[eid])
		}
	}

	if kana.IsLatin(cleanWord(s)) {
		if entryIDs := d.furigana.Get(kana.RomajiToHiragana(cleanWord(s))); entryIDs != nil {
			for _, eid := range entryIDs {
				results = append(results, d.entries[eid])
			}
		}
		if entryIDs := d.furigana.Get(kana.RomajiToKatakana(cleanWord(s))); entryIDs != nil {
			for _, eid := range entryIDs {
				results = append(results, d.entries[eid])
			}
		}
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
		if r := d.english.Get(cw); r != nil {
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
