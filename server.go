package main

import (
	"log"

	"net"

	"fmt"

	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

var healthyairSQLport, healthyairSQLuser, healthyairSQLpassword string
var dbConn *sql.DB
var listener net.Listener

func requestHandler(ctx *fasthttp.RequestCtx) {
	var conn net.Conn
	var err error

	switch string(ctx.Path()[:]) {
	case "/":
		fmt.Println("index()")
		break

	case "/register":
		ctx.SendFile("public/views/register.html")
		break

	case "/register/enter":
		conn, err = listener.Accept()
		if err != nil {
			log.Fatal("Cannot accept connection: ", err)
		}
		registerEnter(ctx, conn)
		conn.Close()
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
	var err error
	healthyairSQLport = os.Getenv("HEALTHYAIR_SQL_PORT")
	if healthyairSQLport == "" {
		log.Fatal("@HEALTHYAIR_SQL_PORT IS NOT SET")
	}
	healthyairSQLuser = os.Getenv("HEALTHYAIR_SQL_USER")
	if healthyairSQLport == "" {
		log.Fatal("@HEALTHYAIR_SQL_USER IS NOT SET")
	}
	healthyairSQLpassword = os.Getenv("HEALTHYAIR_SQL_PASSWORD")
	if healthyairSQLport == "" {
		log.Fatal("@HEALTHYAIR_SQL_PASSWORD IS NOT SET")
	}
	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/healthyair", healthyairSQLuser,
		healthyairSQLpassword))
	if err != nil {
		log.Fatal("@ERR ON OPENING DATABASE", err)
	}
	_, err = dbConn.Exec("SET CHARSET utf8")
	if err != nil {
		log.Fatal("@ERR ON SETTING CHARSET: ", err)
	}
	listener, err = net.Listen("tcp", ":"+healthyairSQLport)
	if err != nil {
		log.Fatal("@ERR ON STARTING LISTENING PORT: ", err)
	}

	err = fasthttp.ListenAndServe("0.0.0.0:80", requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP WEBSERVER ", err)
	}
}
