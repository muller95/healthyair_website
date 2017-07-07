package main

import (
	"log"

	"fmt"

	"html/template"
	"strings"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

type RegisterCard struct {
	Password     string
	Register     string
	Registration string
	Cancel       string
	Name         string
}

func PasswordIsValid(str string) bool {
	return (len(str) >= 7) && (strings.ToLower(str) != str) && strings.ContainsAny(str, "0 & 1 & 2 & 3 & 4 & 5 & 6 & 7 & 8 & 9")
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

	if !govalidator.IsEmail(email) {
		log.Println("@ERR ON WRONG EMAIL SYNTAX")
		return
	}
	if !PasswordIsValid(password) {
		log.Println("@ERR ON WRONG PASSWORD SYNTAX")
		return
	}
	query := "SELECT * FROM users WHERE email=?;"
	//fmt.Println(query)
	rows, err := dbConn.Query(query, email)
	defer rows.Close()
	if err != nil {
		log.Println("@ERR ON QUERY: 'SELECT * FROM users WHERE email=...':", err)
		return
	}
	if rows.Next() {
		log.Println("@ERR ON WRONG EMAIL(ALREADY EXIST)")
		return
	}

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
