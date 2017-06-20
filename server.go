package main

import (
	"log"

	"fmt"

	"github.com/valyala/fasthttp"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()[:]) {
	case "/":
		fmt.Println("index()")
		break

	case "/register":
		fmt.Println("register()")
		break

	case "/register/enter":
		registerEnter(ctx)
		break

	case "/authorization":
		fmt.Println("authorization()")
		break

	default:
		fmt.Println("bad_gateway()")
	}
}
func main() {
	err := fasthttp.ListenAndServe("127.0.0.1:80", requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP WEBSERVER ", err)
	}
}
