package helpers

import (
	"fmt"
	"github.com/gojp/kana"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"strings"
)

func Search(query string) (hits [][]byte) {
	api.Domain = "localhost"

	kana := kana.NewKana()

	isLatin := kana.IsLatin(query)
	isKana := kana.IsKana(query)

	// convert to hiragana and katakana
	romaji := kana.KanaToRomaji(query)

	// handle different types of input differently:
	matches := []string{}
	if isKana {
		// add boost for exact-matching kana
		matches = append(matches, fmt.Sprintf(`
		{"match" :
			{
				"furigana" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 5.0
				}
			}
		}`, query))

		// also look for romaji version in case
		matches = append(matches, fmt.Sprintf(`
		{"match" :
			{
				"romaji" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 2.0
				}
			}
		}`, romaji))
	}
	if !isLatin {
		matches = append(matches, fmt.Sprintf(`
		{"match" :
			{
				"japanese" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 10.0
				}
			}
		}`, query))
	} else {
		// add romaji search term
		matches = append(matches, fmt.Sprintf(`
		{"match" :
			{
				"romaji" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 3.0
				}
			}
		}`, query))

		// add english search term
		matches = append(matches, fmt.Sprintf(`
		{"match" :
			{
				"english" : {
					"query" : "%s",
					"type" : "phrase",
					"boost": 5.0
				}
			}
		}`, query))
	}

	searchJson := fmt.Sprintf(`
		{"query":
			{"bool":
				{
				"should":
					[` + strings.Join(matches, ",") + `],
				"minimum_should_match" : 0,
				"boost": 2.0
				}
			}
		}`)

	out, err := core.SearchRequest(true, "edict", "entry", searchJson, "", 0)
	if err != nil {
		log.Println(err)
	}

	for _, hit := range out.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	return hits
}
