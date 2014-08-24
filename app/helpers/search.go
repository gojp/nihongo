package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gojp/kana"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/revel/revel"
)

func initElasticConnection() {
	elasticURL, _ := revel.Config.String("elastic.url")
	api.Domain = elasticURL

	elasticPort, found := revel.Config.String("elastic.port")
	if found {
		api.Port = string(elasticPort)
	}
}

func executeSearch(searchJson string) (hits [][]byte, err error) {
	out, err := core.SearchRequest("edict", "entry", nil, searchJson)
	if err != nil {
		return hits, err
	}
	for _, hit := range out.Hits.Hits {
		var h interface{}
		h, err = json.Marshal(&hit.Source)
		if err != nil {
			log.Println(err)
		}

		hits = append(hits, h.([]byte))
	}
	return hits, nil
}

func Search(query string) (hits [][]byte, err error) {
	initElasticConnection()

	query = strings.Replace(query, "\"", "\\\"", -1)

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

	hits, err = executeSearch(searchJson)
	return hits, err
}

// ExactSearchVerb is an exact match search on the hiragana and kanji fields
// with a verb-only filter applied.
func ExactSearchVerb(query string, prefix string) (hits [][]byte, err error) {
	initElasticConnection()

	query = strings.Replace(query, "\"", "\\\"", -1)

	searchJson := fmt.Sprintf(`
	{
        "query" : {
            "filtered" : {
                "filter" : {
                    "bool" : {
                        "should" : [
                            { "term" : { "japanese.exact" : "%s" } },
                            { "term" : { "furigana.exact" : "%s" } }
                        ],
                        "must" : [
                            { "prefix" : { "pos" : "%s"}}
                        ]
                    }
                }
            }
        }
	}`, query, query, prefix)
	hits, err = executeSearch(searchJson)
	return hits, err
}

// FuzzySearch returns words similar to the search terms
// provided, and not just exact matches.
func FuzzySearch(query string) (hits [][]byte, err error) {
	initElasticConnection()

	query = strings.Replace(query, "\"", "\\\"", -1)

	searchJson := fmt.Sprintf(`
		{"query":
			{"fuzzy_like_this":
				{
				"fields" : ["romaji", "english"],
				"like_text" : "%s",
				"max_query_terms" : 12
				}
			}
		}`, query)

	hits, err = executeSearch(searchJson)
	return hits, err
}
