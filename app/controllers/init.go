package controllers

import (
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"html/template"
	"strings"
)

func init() {
	revel.TemplateFuncs["add"] = func(a, b int) int { return a + b }
	revel.TemplateFuncs["get"] = func(a []string, b int) string { return a[b] }
	revel.TemplateFuncs["html"] = func(s string) template.HTML { return template.HTML(s) }
	revel.TemplateFuncs["clean"] = func(s string) string {
		s = strings.Replace(s, `<strong>`, ``, -1)
		s = strings.Replace(s, `</strong>`, ``, -1)
		return s
	}
	revel.TemplateFuncs["contains"] = func(a string, b []string) bool {
		for i := range b {
			if a == b[i] {
				return true
			}
		}
		return false
	}

	revel.TemplateFuncs["domain"] = func() string {
		domain, found := revel.Config.String("domain")
		if found {
			return domain
		}
		return ""
	}
	revmgo.ControllerInit()
}
