package main

import (
	"html/template"
	"log"

	"github.com/valyala/fasthttp"
)

type MainPage struct {
	MainPage string
	Main     string
	Buy      string
	AboutUs  string
	Contacts string
}

func mainPage(ctx *fasthttp.RequestCtx, language string) {
	resources := languageResources["en"]
	if language == "ru" {
		resources = languageResources["ru"]
	}

	page := MainPage{MainPage: resources["MainPage"], Main: resources["Main"], Buy: resources["Buy"],
		AboutUs: resources["AboutUs"], Contacts: resources["Contacts"]}
	registerTemplate, err := template.ParseFiles("public/views/main_page.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
	}
	ctx.SetContentType("text/html")
	err = registerTemplate.Execute(ctx, page)
	if err != nil {
		log.Println("Err on executing main page template: ", err)
	}
}
