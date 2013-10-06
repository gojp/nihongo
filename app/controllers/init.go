package controllers

import (
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
)

func init() {
	revel.TemplateFuncs["add"] = func(a, b int) int { return a + b }
	revmgo.ControllerInit()
}
