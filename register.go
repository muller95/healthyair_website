package main

import (
	"fmt"
	"log"
	"strings"

	"unicode/utf8"

	"encoding/json"

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

func PasswordIsStrong(str string) bool {
	return (utf8.RuneCountInString(str) >= 7) && (strings.ToLower(str) != str) &&
		strings.ContainsAny(str, "0 & 1 & 2 & 3 & 4 & 5 & 6 & 7 & 8 & 9")
}

func Register(ctx *fasthttp.RequestCtx, session *Session) {
	email := string(ctx.PostArgs().Peek("email"))
	password := string(ctx.PostArgs().Peek("password"))
	name := string(ctx.PostArgs().Peek("name"))

	resources := languageResources[session.PreferredLanguage]
	resultMap := make(map[string]string)
	errored := false

	if utf8.RuneCountInString(email) == 0 {
		errored = true
		resultMap["email_result"] = "err"
		resultMap["email_message"] = resources["EmptyEmail"]
	} else if !govalidator.IsEmail(email) {
		errored = true
		resultMap["email_result"] = "err"
		resultMap["email_message"] = resources["InvalidEmail"]
	} else {
		resultMap["email_result"] = "ok"
	}

	if utf8.RuneCountInString(name) == 0 {
		errored = true
		resultMap["name_result"] = "err"
		resultMap["name_message"] = resources["EmptyName"]
	} else {
		resultMap["name_result"] = "ok"
	}

	if utf8.RuneCountInString(password) == 0 {
		errored = true
		resultMap["password_result"] = "err"
		resultMap["password_message"] = resources["EmptyPassword"]
	} else if !PasswordIsStrong(password) {
		errored = true
		resultMap["password_result"] = "err"
		resultMap["password_message"] = resources["WeakPassword"]
	} else {
		resultMap["password_result"] = "ok"
	}

	query := "SELECT * FROM users WHERE email=?;"
	//fmt.Println(query)
	rows, err := dbConn.Query(query, email)
	defer rows.Close()
	if err != nil {
		log.Println("@ERR ON QUERY: 'SELECT * FROM users WHERE email=...':", err)
		ctx.Response.SetStatusCode(int(InternalServerError))
		return
	}

	if rows.Next() {
		errored = true
		resultMap["email_result"] = "err"
		resultMap["email_message"] = resources["EmailExists"]
	}

	if errored {
		resultMap["result"] = "err"
		data, err := json.Marshal(resultMap)
		if err != nil {
			log.Println("Err encoding register answer: ", err)
			ctx.Response.SetStatusCode(int(InternalServerError))
			return
		}

		ctx.Response.SetStatusCode(int(Ok))
		ctx.Write(data)
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

	resultMap["result"] = "ok"
	data, err := json.Marshal(resultMap)
	if err != nil {
		log.Println("Err encoding register answer: ", err)
		ctx.Response.SetStatusCode(int(InternalServerError))
		return
	}

	ctx.Response.SetStatusCode(int(Ok))
	ctx.Write(data)
}
