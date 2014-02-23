// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/robfig/revel"


type tApp struct {}
var App tApp


func (_ tApp) Search(
		query string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "query", query)
	return revel.MainRouter.Reverse("App.Search", args).Url
}

func (_ tApp) Details(
		query string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "query", query)
	return revel.MainRouter.Reverse("App.Details", args).Url
}

func (_ tApp) SearchGet(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.SearchGet", args).Url
}

func (_ tApp) About(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.About", args).Url
}

func (_ tApp) SavePhrase(
		phrase string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "phrase", phrase)
	return revel.MainRouter.Reverse("App.SavePhrase", args).Url
}

func (_ tApp) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Register", args).Url
}

func (_ tApp) LoginPage(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.LoginPage", args).Url
}

func (_ tApp) Profile(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Profile", args).Url
}

func (_ tApp) SaveUser(
		user interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	return revel.MainRouter.Reverse("App.SaveUser", args).Url
}

func (_ tApp) Login(
		email string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("App.Login", args).Url
}

func (_ tApp) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Logout", args).Url
}

func (_ tApp) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Index", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


