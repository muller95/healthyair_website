package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

func authorizeEnter(ctx *fasthttp.RequestCtx, session *Session) {
	email := string(ctx.PostArgs().Peek("email"))
	password := string(ctx.PostArgs().Peek("password"))

	query := "SELECT id FROM users WHERE email=\"" + email + "\" && passwd=\"" + password + "\";"
	//fmt.Println(query)
	rows, err := dbConn.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println("@ERR ON QUERY: 'SELECT id FROM users WHERE email=..., passwd=...':", err)
		return
	}
	if rows.Next() {
		var uid uint64
		err = rows.Scan(&uid)
		if err != nil {
			log.Println("@ERR ON SCANNING ROWS: ", err)
			return
		}
		session.UserID = uid
		session.Authorized = true
		SessionSetUserID(session)
	} else {
		log.Println("@ERR ON WRONG EMAIL OR PASSWORD")
		return
	}
}
