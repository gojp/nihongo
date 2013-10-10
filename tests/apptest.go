package tests

import "github.com/robfig/revel"

type AppTest struct {
	revel.TestSuite
}

func (t AppTest) Before() {
	println("Set up")
}

func (t AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatHelloSearchPageWorks() {
	t.Get("/hello")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatKonnichiwaSearchPageWorks() {
	t.Get("/今日は")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) TestThatDoublePlusSearchWorks() {
	t.Get("/今日は++")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t AppTest) After() {
	println("Tear down")
}
