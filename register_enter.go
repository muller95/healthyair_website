package main

import (
	"log"

	"fmt"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

//func send_to_mysql()

func registerEnter(ctx *fasthttp.RequestCtx) {
	var email, passwd string
	var err error

	if ctx.IsPost() {
		pbody := string(ctx.PostBody())
		fmt.Println(string(pbody))
		//request must contain 2 args(email, password)
		if strings.Contains(pbody, "&") {
			email = pbody[strings.Index(pbody, "=")+1 : strings.Index(pbody, "&")]
			pbody = strings.TrimPrefix(pbody, pbody[:strings.Index(pbody, "&")+1])
		}
		passwd = pbody[strings.Index(pbody, "=")+1:]
		//stmt := fmt.Sprintf("INSERT INTO users (email, passwd) VALUES (\"%s\", \"%s\")", email, passwd)
		_, err = dbConn.Prepare("INSERT INTO users (email, passwd) VALUES (\"?\", \"?\")")
		if err != nil {
			log.Fatal("@ERR ON PREPARING STMT: 'INSERT INTO users (email, passwd) VALUES (\"?\", \"?\")':", err)
		}

		//fmt.Println(stmt)
		_, err = dbConn.Exec(email, passwd)
		if err != nil {
			log.Fatal("@ERR ON EXECUTING STMT: 'INSERT INTO users (email, passwd) VALUES (?, ?)'", err)
		}
	} else {
		log.Fatal("@ERR ON METHOD-REQUEST")
	}
}
