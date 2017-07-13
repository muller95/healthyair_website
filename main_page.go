package main

import (
	"html/template"
	"log"

	"bytes"

	"github.com/valyala/fasthttp"
)

type Navbar struct {
	Main         string
	Buy          string
	AboutUs      string
	Contacts     string
	Language     string
	SignIn       string
	Registration string
	Modals       template.HTML
}

type Cards struct {
	Conveniently      string
	ConvenientlyText  string
	Fast              string
	FastText          string
	Qualitatively     string
	QualitativelyText string
}

type MainPage struct {
	MainPage     string
	Navbar       template.HTML
	MainPageText string
	Cards        template.HTML
}

type NavbarModals struct {
	Registration string
	Name         string
	Password     string
	Register     string
	Cancel       string
}

func executeModals(resources map[string]string, session *Session) template.HTML {
	var navbarModals NavbarModals

	navbarModals.Registration = resources["Registration"]
	navbarModals.Name = resources["Name"]
	navbarModals.Password = resources["Password"]
	navbarModals.Register = resources["Register"]
	navbarModals.Cancel = resources["Cancel"]

	navbarModalsTemplate, err := template.ParseFiles("public/views/navbar_modals_template.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
		return ""
	}

	navbarModalsWriter := &bytes.Buffer{}
	err = navbarModalsTemplate.Execute(navbarModalsWriter, navbarModals)
	if err != nil {
		log.Println("Err on executing navbar template: ", err)
		return ""
	}

	return template.HTML(navbarModalsWriter.Bytes())
}

func executeNavbar(resources map[string]string, session *Session) template.HTML {
	var navbar Navbar

	navbar.AboutUs = resources["AboutUs"]
	navbar.Buy = resources["Buy"]
	navbar.Contacts = resources["Contacts"]
	navbar.Main = resources["Main"]
	navbar.Language = session.PreferredLanguage
	navbar.SignIn = resources["SignIn"]
	navbar.Registration = resources["Registration"]
	navbar.Modals = executeModals(resources, session)

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

func executeCards(resources map[string]string) template.HTML {
	var cards Cards

	cards.Conveniently = resources["Conveniently"]
	cards.ConvenientlyText = resources["ConvenientlyText"]
	cards.Fast = resources["Fast"]
	cards.FastText = resources["FastText"]
	cards.Qualitatively = resources["Qualitatively"]
	cards.QualitativelyText = resources["QualitativelyText"]

	cardsTemplate, err := template.ParseFiles("public/views/cards_template.html")
	if err != nil {
		log.Println("Err on parsing main page template: ", err)
		return ""
	}

	cardsWriter := &bytes.Buffer{}
	err = cardsTemplate.Execute(cardsWriter, cards)
	if err != nil {
		log.Println("Err on executing cards template: ", err)
		return ""
	}

	return template.HTML(cardsWriter.Bytes())
}

func mainPage(ctx *fasthttp.RequestCtx, session *Session) {
	var mainPage MainPage

	log.Println(session)
	resources := languageResources[session.PreferredLanguage]

	mainPage.MainPage = resources["MainPage"]
	mainPage.Navbar = executeNavbar(resources, session)
	mainPage.Cards = executeCards(resources)
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
