package tests

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type NihongoSuite struct{}

var _ = Suite(&NihongoSuite{})

func (s *NihongoSuite) TestSearchResults(c *C) {
	// some basic checks
	c.Check(1, Equals, 2)
}
