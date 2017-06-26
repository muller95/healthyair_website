package main

import (
	"log"

	"fmt"

	"os"

	"database/sql"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

var languageResources map[string]map[string]string

var healthyairSQLport, healthyairSQLuser, healthyairSQLpassword string
var dbConn *sql.DB

//var listener net.Listener

func requestHandler(ctx *fasthttp.RequestCtx) {
	//var conn net.Conn
	//var err error
	language := "en"
	acceptLanguages := strings.Split(string(ctx.Request.Header.Peek("Accept-Language")), ";")
	if len(acceptLanguages) > 0 {
		acceptLanguage := strings.Split(acceptLanguages[0], ",")
		if len(acceptLanguage) > 1 {
			language = acceptLanguage[1]
		}
	}

	switch string(ctx.Path()[:]) {
	case "/":
		fmt.Println("index()")
		break

	case "/register":
		registerForm(ctx, language)
		break

	case "/register/enter":
		//fmt.Println("here1")
		/*conn, err = listener.Accept()
		if err != nil {
			log.Fatal("Cannot accept connection: ", err)
		}*/
		registerEnter(ctx)
		//fmt.Println("here3")
		//conn.Close()
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

func initResources() {
	languageResources = make(map[string]map[string]string)

	languageResources["en"] = make(map[string]string)
	languageResources["ru"] = make(map[string]string)

	languageResources["en"]["Password"] = "Password"
	languageResources["ru"]["Password"] = "Пароль"

	languageResources["en"]["Register"] = "Register"
	languageResources["ru"]["Register"] = "Зарегистрироваться"

	languageResources["en"]["Registration"] = "Registration"
	languageResources["ru"]["Registration"] = "Регистрация"

	languageResources["en"]["Cancel"] = "Cancel"
	languageResources["ru"]["Cancel"] = "Отмена"

	languageResources["en"]["Name"] = "Name"
	languageResources["ru"]["Name"] = "Имя"
}

func main() {
	var err error
	/*healthyairSQLport = os.Getenv("HEALTHYAIR_SQL_PORT")
	if healthyairSQLport == "" {
		log.Fatal("@HEALTHYAIR_SQL_PORT IS NOT SET")
	}*/

	initResources()

	healthyairSQLuser = os.Getenv("HEALTHYAIR_SQL_USER")
	if healthyairSQLuser == "" {
		log.Fatal("@HEALTHYAIR_SQL_USER IS NOT SET")
	}
	healthyairSQLpassword = os.Getenv("HEALTHYAIR_SQL_PASSWORD")
	if healthyairSQLpassword == "" {
		log.Fatal("@HEALTHYAIR_SQL_PASSWORD IS NOT SET")
	}
	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/healthyair?charset=utf8", healthyairSQLuser,
		healthyairSQLpassword))
	if err != nil {
		log.Fatal("@ERR ON OPENING DATABASE", err)
	}
	/*_, err = dbConn.Exec("SET CHARSET utf8")
	if err != nil {
		log.Fatal("@ERR ON SETTING CHARSET: ", err)
	}*/
	/*listener, err = net.Listen("tcp", ":"+healthyairSQLport)
	if err != nil {
		log.Fatal("@ERR ON STARTING LISTENING PORT: ", err)
	}*/

	err = fasthttp.ListenAndServe("0.0.0.0:80", requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP WEBSERVER ", err)
	}
	dbConn.Close()
}
