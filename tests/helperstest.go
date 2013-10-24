package tests

import (
	"github.com/gojp/nihongo/app/helpers"
	"github.com/robfig/revel"
)

type HelpersTest struct {
	revel.TestSuite
}

func (s *HelpersTest) Before() {
}

func (s *HelpersTest) TestMakeStrong() {
	// some basic checks
	s.Assert(helpers.MakeStrong("gopher") == "<strong>gopher</strong>")
}

func (s *HelpersTest) After() {
}
