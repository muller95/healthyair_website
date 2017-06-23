package main

import (
	"log"

	"fmt"

	"strings"

	"net"

	//"database/sql"

	"github.com/valyala/fasthttp"
	//_ "github.com/go-sql-driver/mysql"
)

//func send_to_mysql()

func registerEnter(ctx *fasthttp.RequestCtx, conn net.Conn) {
	var email, passwd string
	var err error

	if ctx.IsPost() {
		pbody := string(ctx.PostBody())
		fmt.Println(string(pbody))
		//request must contains 2 args(email, password)
		if strings.Contains(pbody, "&") {
			email = pbody[strings.Index(pbody, "=")+1 : strings.Index(pbody, "&")]
			pbody = strings.TrimPrefix(pbody, pbody[:strings.Index(pbody, "&")+1])
		}
		passwd = pbody[strings.Index(pbody, "=")+1:]

		_, err = dbConn.Prepare("INSERT INTO users (email, passwd) VALUES (?, ?)")
		if err != nil {
			log.Fatal("@ERR ON PREPARING STMT: 'INSERT INTO users (email, passwd) VALUES (?, ?)'")
		}
		_, err = dbConn.Exec(email, passwd)
		if err != nil {
			log.Fatal("@ERR ON EXECUTING STMT: 'INSERT INTO users (email, passwd) VALUES (?, ?)'")
		}
	} else {
		log.Fatal("@ERR ON METHOD-REQUEST")
	}
}
