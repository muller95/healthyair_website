package main

import (
	"html/template"
	"log"

	"bytes"

	"github.com/valyala/fasthttp"
)

type Navbar struct {
	Main     string
	Buy      string
	AboutUs  string
	Contacts string
}

type MainPage struct {
	MainPage     string
	Navbar       template.HTML
	MainPageText string
}

func executeNavbar(resources map[string]string) template.HTML {
	var navbar Navbar

	navbar.AboutUs = resources["AboutUs"]
	navbar.Buy = resources["Buy"]
	navbar.Contacts = resources["Contacts"]
	navbar.Main = resources["Main"]

	navbarTemplate, err := template.ParseFiles("public/views/nav_template.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
		return ""
	}

	navbarWriter := &bytes.Buffer{}
	err = navbarTemplate.Execute(navbarWriter, navbar)
	if err != nil {
		log.Println("Err on executing navbar template: ", err)
		return ""
	}

	return template.HTML(navbarWriter.Bytes())
}

func mainPage(ctx *fasthttp.RequestCtx, language string) {
	var mainPage MainPage

	resources := languageResources["en"]
	if language == "ru" {
		resources = languageResources["ru"]
	}

	mainPage.MainPage = resources["MainPage"]
	mainPage.Navbar = executeNavbar(resources)
	mainPage.MainPageText = resources["MainPageText"]

	registerTemplate, err := template.ParseFiles("public/views/main_page.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
	}
	ctx.SetContentType("text/html")
	err = registerTemplate.Execute(ctx, mainPage)
	if err != nil {
		log.Println("Err on executing main page template: ", err)
	}
}
