package main

import (
	"log"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"

	"html/template"
)

//func send_to_mysql()

type RegisterCard struct {
	Password     string
	Register     string
	Registration string
	Cancel       string
	Name         string
}

func registerForm(ctx *fasthttp.RequestCtx, language string) {
	resources := languageResources["en"]
	if language == "ru" {
		resources = languageResources["ru"]
	}

	card := RegisterCard{Password: resources["Password"], Register: resources["Register"],
		Cancel: resources["Cancel"], Name: resources["Name"], Registration: resources["Registration"]}
	registerTemplate, err := template.ParseFiles("public/views/register.html")
	if err != nil {
		log.Println("Err on parsing register card template: ", err)
	}
	ctx.SetContentType("text/html")
	err = registerTemplate.Execute(ctx, card)
	if err != nil {
		log.Println("Err on executing register template: ", err)
	}
}

func registerEnter(ctx *fasthttp.RequestCtx) {
	// var err error

	email := string(ctx.PostArgs().Peek("email"))
	password := string(ctx.PostArgs().Peek("password"))
	name := string(ctx.PostArgs().Peek("name"))
	//stmt := fmt.Sprintf("INSERT INTO users (email, passwd) VALUES (\"%s\", \"%s\")", email, passwd)
	stmt, err := dbConn.Prepare("INSERT INTO users (email, passwd, name) VALUES (?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		log.Println("@ERR ON PREPARING STMT: 'INSERT INTO users (email, passwd, name) VALUES (?, ?, ?)':", err)
		return
	}

	fmt.Println(stmt)
	_, err = stmt.Exec(email, password, name)
	// _, err = dbConn.Exec(email, password, name)
	if err != nil {
		log.Println("@ERR ON EXECUTING STMT: INSERT INTO users (email, passwd, name) VALUES ('?', '?', '?')", err)
		return
	}

}
