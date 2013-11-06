package controllers

import (
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"html/template"
)

func init() {
	revel.TemplateFuncs["add"] = func(a, b int) int { return a + b }
	revel.TemplateFuncs["get"] = func(a []string, b int) string { return a[b] }
	revel.TemplateFuncs["html"] = func(s string) template.HTML { return template.HTML(s) }
	revel.TemplateFuncs["contains"] = func(a string, b []string) bool {
		for i := range b {
			if a == b[i] {
				return true
			}
		}
		return false
	}
	revmgo.ControllerInit()
}
