package tests

import (
	"github.com/gojp/nihongo/app/helpers"
	"github.com/revel/revel"
)

type HelpersTest struct {
	revel.TestSuite
}

func (s *HelpersTest) Before() {
}

func (s *HelpersTest) After() {
}

func (s *HelpersTest) TestMakeStrong() {
	// some basic checks
	s.Assert(helpers.MakeStrong("gopher") == "<strong>gopher</strong>")
}

func (s *HelpersTest) TestConvertQueryToKana() {
	// some basic checks
	h, k := helpers.ConvertQueryToKana("tesuto")
	s.Assert(h == "てすと")
	s.Assert(k == "テスト")
}
