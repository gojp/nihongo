package controllers

import "github.com/robfig/revel"
import "html/template"

func init() {
	revel.TemplateFuncs["add"] = func(a, b int) int { return a + b }
	revel.TemplateFuncs["html"] = func(s string) template.HTML { return template.HTML(s) }
}
