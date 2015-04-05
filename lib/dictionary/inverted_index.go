package dictionary

import (
	"bytes"
	"encoding/binary"
	"log"
	"sort"

	"github.com/gojp/nihongo/lib/bloomfilter"
)

type InvertedIndex struct {
	MaxReferences int // maximum number of references stored per key
	entries       map[string]*IndexEntry
}

func NewInvertedIndex(maxRef int) *InvertedIndex {
	return &InvertedIndex{
		MaxReferences: maxRef,
		entries:       map[string]*IndexEntry{},
	}
}

type IndexEntry struct {
	references []Reference

	// a bloom filter to probabilistically store references that go beyond MaxReferences
	filter bloomfilter.BloomFilter
}

func newIndexEntry(filterSize uint) *IndexEntry {
	return &IndexEntry{
		references: []Reference{},
		filter:     *bloomfilter.NewEstimated(filterSize, 0.01),
	}
}

type Reference struct {
	id    EntryID
	score float64
}

// ByScore sorts references in descending order of score.
type ByScore []Reference

// Len is part of sort.Interface.
func (s ByScore) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s ByScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface.
func (s ByScore) Less(i, j int) bool {
	return s[i].score >= s[j].score
}

// Insert inserts a reference to the index
func (i *InvertedIndex) Insert(key string, id EntryID, score float64) {
	var (
		ie *IndexEntry
		ok bool
	)
	if ie, ok = i.entries[key]; !ok {
		ie = newIndexEntry(1000)
		i.entries[key] = ie
	}
	r := Reference{
		id:    id,
		score: score,
	}
	ie.references = append(ie.references, r)
	sort.Sort(ByScore(ie.references)) // todo: use BST instead
	if len(ie.references) > i.MaxReferences {
		// throw away the last entry that exceeds the limit, but add it to the
		// bloom filter
		extra := ie.references[i.MaxReferences]
		ie.references = ie.references[0:i.MaxReferences]

		ie.filter.Add(getIDBytes(extra.id))
	}
}

// Get fetches a slice of references for the given key, or returns a nil
// slice of no slice is found.
func (i *InvertedIndex) Get(key string) []Reference {
	e, ok := i.entries[key]
	if !ok {
		return nil
	}
	return e.references
}

// Test returns whether the given id might be contained at the key.
// It does this probabilistically, using a bloom filter.
func (i *InvertedIndex) Test(key string, id EntryID) bool {
	e, ok := i.entries[key]
	if !ok {
		return false
	}
	return e.filter.Test(getIDBytes(id))
}

func getIDBytes(id EntryID) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, id)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
