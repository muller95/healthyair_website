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
		ctx.SendFile("public/views/register.html")
		break

	case "/register/enter":
		registerEnter(ctx)
		break

	case "/authorization":
		fmt.Println("authorization()")
		break

	default:
		fasthttp.FSHandler("public/", 0)(ctx)

		// default:
		// fmt.Println("bad_gateway()")
	}
}
func main() {
	err := fasthttp.ListenAndServe("0.0.0.0:80", requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP WEBSERVER ", err)
	}
}
