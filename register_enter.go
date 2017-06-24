package main

import (
	"log"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

//func send_to_mysql()

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
