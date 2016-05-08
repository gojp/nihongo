package dictionary

import "testing"

var getTests = []string{
	"しけん",
	"てつだう",
	"手伝う",
	"ふつう",
	"普通",
	"ふつ",
}

var notIncluded = []string{
	"ふ",
	"ふつつ",
	"ふつうう",
}

func TestGet(t *testing.T) {
	r := NewRadixTree()
	for i, entry := range getTests {
		r.Insert(entry, EntryID(i))
		got := r.Get(entry)
		if len(got) != 1 {
			t.Fatalf("%q len(got) = %d, want %d", entry, len(got), 1)
		}
		if got[0] != EntryID(i) {
			t.Fatalf("got[0] = %q, want %q", got[0], entry)
		}
	}

	for _, entry := range notIncluded {
		got := r.Get(entry)
		if len(got) != 0 {
			t.Fatalf("%q len(got) = %d, want %d", entry, len(got), 0)
		}
	}
}

func TestFindWordsWithPrefix(t *testing.T) {
	r := NewRadixTree()
	for i, entry := range getTests {
		r.Insert(entry, EntryID(i))
	}

	got := len(r.FindWordsWithPrefix("普通", 10))
	if got != 1 {
		t.Fatalf("%q len(got) = %d, want %d", "普通", got, 1)
	}

	got = len(r.FindWordsWithPrefix("ふ", 10))
	if got != 2 {
		t.Fatalf("%q len(got) = %d, want %d", "ふ", got, 2)
	}
}
