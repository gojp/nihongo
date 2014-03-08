package helpers

import (
	"strconv"

	"github.com/gojp/kanjidic2"
	"github.com/gojp/radicals"
)

type scToKanji map[int][]string
type radMap map[string]scToKanji

func LoadRadMap() (radMap, error) {
	radMap := make(radMap)
	r, err := radicals.ParseRadkfile("radkfile.utf-8")
	if err != nil {
		return radMap, err
	}
	k, err := kanjidic2.ParseKanjiDic2("kanjidic2.xml")

	if err != nil {
		return radMap, err
	}

	for radStr, rad := range r {
		key := radStr + "_" + strconv.Itoa(rad.StrokeCount)
		radMap[key] = scToKanji{}
		for _, l := range rad.Kanji {
			ksc := k[l].StrokeCount
			radMap[key][ksc] = append(radMap[key][ksc], l)
		}
	}
	return radMap, nil
}
