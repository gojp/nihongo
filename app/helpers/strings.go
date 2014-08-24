package helpers

import (
	"github.com/gojp/japanese"
	"github.com/gojp/kana"
)

var verbForms = map[int]string{
	japanese.TeForm:             "te",
	japanese.Infinitive:         "infinitive",
	japanese.PresentIndicative:  "present indicative",
	japanese.Presumptive:        "presumptive",
	japanese.Imperative:         "imperative",
	japanese.PastIndicative:     "past indicative",
	japanese.PastPresumptive:    "past presumptive",
	japanese.PresentProgressive: "present progressive",
	japanese.PastProgressive:    "past progressive",
	japanese.Provisional:        "provisional",
	japanese.Conditional:        "conditional",
	japanese.Potential:          "potential",
	japanese.Causative:          "causative",
}

// MakeStrong wraps a string in strong tags.
func MakeStrong(query string) string {
	return "<strong>" + query + "</strong>"
}

// ConvertQueryToKana converts the query to hiragana and katakana.
// If it's already in hiragana or katakana, it will just be the same.
func ConvertQueryToKana(query string) (hiragana, katakana string) {
	h := kana.RomajiToHiragana(query)
	k := kana.RomajiToKatakana(query)
	return h, k
}

// IdentifyVerb returns whether or not a given string can be identified
// as being in a verb form, and the associate form if applicable.
func IdentifyVerb(word string) (isVerb bool, form string) {
	f := japanese.IdentifyForm(word)
	if form, isVerb = verbForms[f]; isVerb {
		return isVerb, form
	}
	return false, ""
}
