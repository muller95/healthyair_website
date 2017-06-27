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

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDVHqu2YJUDLjYdCDMmRlHJX1KI
eiuiU6i+JbFrU/AylGXt2oCaGxGMiNqu1UIhxaG1Z6sozmgEFIZ9PPScAXghPm54
eKspTmQ3oFtWGcyq9ury0HEezYDi7TZHlv2wKIVDacBHivNgsVQIhuN3ICHFFHMq
9O0aN2dZFXB0rOImgQIDAQAB
-----END PUBLIC KEY-----
`)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDVHqu2YJUDLjYdCDMmRlHJX1KIeiuiU6i+JbFrU/AylGXt2oCa
GxGMiNqu1UIhxaG1Z6sozmgEFIZ9PPScAXghPm54eKspTmQ3oFtWGcyq9ury0HEe
zYDi7TZHlv2wKIVDacBHivNgsVQIhuN3ICHFFHMq9O0aN2dZFXB0rOImgQIDAQAB
AoGAISMJs+vEf6AZzd3Ohi783ICzxoCodC7p19bohTWh7VthleAZityWl/FXf0Ot
aq7d++TImimtxqSiXKqzpeYclVXkprwZfppCVUfK0YBf/JdBBZM0kaq1d+iKjTo+
jPUy9U/4Tqd5o8oXWG8lS9vRiEUKO4VixEJsbV+3yaeC6AECQQD5A2VAarih46Xd
aNmY9nC6J5MX0te9/IxODT+9KhDhY1JS9A6FUbm/sYl8NfLzsTBv2mD0pE54gXpF
wY3GnnehAkEA2xl1+7NIiLQ7sCcpRFqzP0kBdIg/ir2Ahq3imV4RcQoP93RHGJ3m
9pSdhz4lvA5mJtH0RFwMe2a0Wa2PpBPC4QJAVMa7KfsrcLI4PfD8Y/9C0Z23jlzR
5nScr9YC5Tv1E0blOCiu6OSyAHlI/WjAlga1Ht+SMrfdn1k1b5o90mkRAQJBAIHq
GvdgW0YT+MB+uA176oU/+MjscSEHNMqnGJHwIU9xs/36yJ1kI6tae/3Rb/aOYyvp
mnleS1hwkcgLDf0waoECQEzihTwr0Mdg7swdcMejvPlSB32kr8JEfd3z71XTdwrz
shapI5224UYfNtBAzSS44jbeM3lBjhWpQvFl3csCddM=
-----END RSA PRIVATE KEY-----
`)

//var listener net.Listener

func requestHandler(ctx *fasthttp.RequestCtx) {
	//var conn net.Conn
	//var err error
	log.Println("request")
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

	// err = fasthttp.ListenAndServe("0.0.0.0:80", requestHandler)
	err = fasthttp.ListenAndServeTLS(":8080", "ssl_cert/server.crt", "ssl_cert/server.key",
		requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP WEBSERVER ", err)
	}
	dbConn.Close()
}
