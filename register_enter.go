package main

import (
	"log"

	"fmt"

	"strings"

	"github.com/valyala/fasthttp"
)

//func send_to_mysql()

func registerEnter(ctx *fasthttp.RequestCtx) {
	if ctx.IsPost() {
		pbody := string(ctx.PostBody())
		fmt.Println(string(pbody))
		//request must contains 5 args(name, email, login, password, confirm_password)
		for i := 0; i < 4; i++ { //parsing first 4 args...
			if strings.Contains(pbody, "&") {
				name := pbody[:strings.Index(pbody, "=")]
				filling := pbody[strings.Index(pbody, "=")+1 : strings.Index(pbody, "&")]
				//HERE need to send info to mysql
				pbody = strings.TrimPrefix(pbody, pbody[:strings.Index(pbody, "&")+1])
				fmt.Println("name = " + name + " | " + "filling = " + filling)
			}
		}
		//parsing last arg
		name := pbody[:strings.Index(pbody, "=")]
		filling := pbody[strings.Index(pbody, "=")+1:]
		//HERE need to send info to mysql

		fmt.Println("name = " + name + " | " + "filling = " + filling)
	} else {
		log.Fatal("@ERR ON METHOD-REQUEST")
	}
}
