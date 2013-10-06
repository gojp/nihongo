package controllers

import (
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"html/template"
)

func init() {
	revel.TemplateFuncs["add"] = func(a, b int) int { return a + b }
	revel.TemplateFuncs["html"] = func(s string) template.HTML { return template.HTML(s) }
	revmgo.ControllerInit()
}
