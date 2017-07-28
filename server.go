package main

import (
	"log"
	"time"

	"fmt"

	"os"

	"database/sql"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	tarantool "github.com/tarantool/go-tarantool"
	"github.com/valyala/fasthttp"
)

var languageResources map[string]map[string]string

var healthyairSQLport, healthyairSQLuser, healthyairSQLpassword string
var dbConn *sql.DB

var healthyairTARANTOOLserver string
var healthyairTARANTOOLopts tarantool.Opts

var healthyairTARANTOOLclient *tarantool.Connection

type RestCode uint32

const (
	Ok                  RestCode = 200
	NotFound            RestCode = 404
	SessionExpired      RestCode = 471
	InternalServerError RestCode = 500
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	var session Session
	var rc RestCode
	var c fasthttp.Cookie

	language := "en"
	acceptLanguages := strings.Split(string(ctx.Request.Header.Peek("Accept-Language")), ";")
	if len(acceptLanguages) > 0 {
		acceptLanguage := strings.Split(acceptLanguages[0], ",")
		if len(acceptLanguage) > 1 {
			language = acceptLanguage[1]
		}
	}

	if len(ctx.Request.Header.Cookie("session_id")) == 0 {
		session, rc = SessionStart(language)
		if rc != Ok {
			ctx.Response.SetStatusCode(int(rc))
			return
		}

		c.SetKey("session_id")
		c.SetValue(session.SessionID)
		ctx.Response.Header.SetCookie(&c)
	} else {
		session, rc = SessionGet(string(ctx.Request.Header.Cookie("session_id")))
		if rc == NotFound || rc == SessionExpired {
			session, rc = SessionStart(language)
			if rc != Ok {
				ctx.Response.SetStatusCode(int(rc))
				return
			}

			c.SetKey("session_id")
			c.SetValue(session.SessionID)
			ctx.Response.Header.SetCookie(&c)
		} else if rc != Ok {
			ctx.Response.SetStatusCode(int(rc))
			return
		}
	}

	switch string(ctx.Path()[:]) {
	case "/":
		mainPage(ctx, &session)
		break

	case "/register":
		Register(ctx, &session)
		break

	case "/authorize":
		fmt.Println("authorize()")
		break

	case "/authorize_enter":
		authorizeEnter(ctx, &session)
		break

	case "/set_preferred_language":
		preferedLanguage := string(ctx.PostArgs().Peek("language"))
		if preferedLanguage != "ru" {
			preferedLanguage = "en"
		}

		session.PreferredLanguage = preferedLanguage
		rc = SessionSetPreferredLanguage(&session)

		ctx.Response.SetStatusCode(int(rc))
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

	languageResources["en"]["Register"] = "Register"
	languageResources["ru"]["Register"] = "Зарегистрироваться"

	languageResources["en"]["Password"] = "Password"
	languageResources["ru"]["Password"] = "Пароль"

	languageResources["en"]["Cancel"] = "Cancel"
	languageResources["ru"]["Cancel"] = "Отмена"

	languageResources["en"]["Name"] = "Name"
	languageResources["ru"]["Name"] = "Имя"

	languageResources["en"]["MainPage"] = "Main page"
	languageResources["ru"]["MainPage"] = "Главная страница"

	languageResources["en"]["Main"] = "Main"
	languageResources["ru"]["Main"] = "Главная"

	languageResources["en"]["Buy"] = "Buy"
	languageResources["ru"]["Buy"] = "Купить"

	languageResources["en"]["AboutUs"] = "About us"
	languageResources["ru"]["AboutUs"] = "О нас"

	languageResources["en"]["Contacts"] = "Contacts"
	languageResources["ru"]["Contacts"] = "Контакты"

	languageResources["en"]["MainPageText"] = `
	Healthy Air project provides you a new room weather station. It measures humidity,
	temperature and carbon dioxide concentration in your room. Its uniqueness lies in the fact that it
	collects all data on the webserver and you can monitor data from several meteostations,
	also it gives you some advices to make air in the room better.  Bereathe easily with Healthy Air!
	`
	languageResources["ru"]["MainPageText"] = `
	Проект Healthy Air представляет вам новую комнатную метеостанцию. Она измеряет влажность, 
	температуру и концентрацию углекислого газа в вашей комнате. Её уникальность состоит в том,
	что она собирает все данные на сервере, и вы можете отслеживать данные с нескольких метеостанций,
	также она даёт некоторые советы, как сделать воздух в комнате лучше. 
	Дышите легко, вместе с Healthy Air!
	`

	languageResources["en"]["Conveniently"] = "CONVENIENTLY"
	languageResources["ru"]["Conveniently"] = "УДОБНО"

	languageResources["en"]["ConvenientlyText"] =
		`
	Learn about the need to ventilate the room with your mobile device.
	`
	languageResources["ru"]["ConvenientlyText"] = `
	Узнавайте о необходимости проветрить помещение с помощью Вашего мобильного устройства.
	`

	languageResources["en"]["Fast"] = "FAST"
	languageResources["ru"]["Fast"] = "БЫСТРО"

	languageResources["en"]["FastText"] = `
	There is no need to take readings directly from the weather station, you can monitor them using 
	the web-site.
	`
	languageResources["ru"]["FastText"] = `
	Нет необходимости снимать показания непосредственно с метеостанции, можно следить за ними
	при помощи веб-сайта.
	`

	languageResources["en"]["Qualitatively"] = "QUALITATIVELY"
	languageResources["ru"]["Qualitatively"] = "КАЧЕСТВЕННО"

	languageResources["en"]["QualitativelyText"] = `
	The weather station will be your reliable assistant in maintaining the cleanliness of the house.
	`
	languageResources["ru"]["QualitativelyText"] = `
	Метостанция станет вашим надёжным помощником в поддержании чистоты дома.
	`

	languageResources["en"]["SignIn"] = "Sign in"
	languageResources["ru"]["SignIn"] = "Войти"

	languageResources["en"]["Registration"] = "Registration"
	languageResources["ru"]["Registration"] = "Регистрация"

	languageResources["en"]["RegisterOK"] = "You have registered successfully."
	languageResources["ru"]["RegisterOK"] = "Вы успешно зарегистрировались."

	languageResources["en"]["RegisterFail"] = "Some error occured, try register again later."
	languageResources["ru"]["RegisterFail"] = "Произошла какая-то ошибка," +
		"попробуйте зарегистрироваться позже."

	languageResources["en"]["EmptyEmail"] = "Enter your email, please."
	languageResources["ru"]["EmptyEmail"] = "Введите свой email, пожалуйста."

	languageResources["en"]["InvalidEmail"] = "Enter valid email, please."
	languageResources["ru"]["InvalidEmail"] = "Введите корректный email, пожалуйста."

	languageResources["en"]["EmptyName"] = "Enter your name, please."
	languageResources["ru"]["EmptyName"] = "Ввеедите ваше имя, пожалуйста."

	languageResources["en"]["EmptyPassword"] = "Enter password, please."
	languageResources["ru"]["EmptyPassword"] = "Введите пароль, пожалуйста."

	languageResources["en"]["WeakPassword"] = "Your password is too weak."
	languageResources["ru"]["WeakPassword"] = "Ваш пароль слишком слабый."

	languageResources["en"]["EmailExists"] = "This email alredy registered."
	languageResources["ru"]["EmailExists"] = "Этот email уже зарегистрирован."

	languageResources["en"]["RegisterHint"] = "The password must contains from seven symbols" +
		", and contains at least one capital letter and one digit."
	languageResources["ru"]["RegisterHint"] = "Пароль должен состоять как минимум из семи букв," +
		", содержать как минимум одну цифру и одну заглавную букву."

	languageResources["en"]["Authorization"] = "Authorization"
	languageResources["ru"]["Authorization"] = "Авторизация"
}

func main() {
	var err error
	initResources()

	healthyairSQLuser = os.Getenv("HEALTHYAIR_SQL_USER")
	if healthyairSQLuser == "" {
		log.Fatal("Err: HEALTHYAIR_SQL_USER is not set")
	}
	healthyairSQLpassword = os.Getenv("HEALTHYAIR_SQL_PASSWORD")
	if healthyairSQLpassword == "" {
		log.Fatal("Err: HEALTHYAIR_SQL_PASSWORD is not set")
	}

	certificatePath := os.Getenv("HEALTHYAIR_CERTIFICATE_PATH")
	if certificatePath == "" {
		log.Fatal("Err: HEALTHYAIR_CERTIFICATE_PATH is not set")
	}
	keyPath := os.Getenv("HEALTHYAIR_KEY_PATH")
	if keyPath == "" {
		log.Fatal("Err: HEALTHYAIR_KEY_PATH is not set")
	}

	serverPort := os.Getenv("HEALTHYAIR_SERVER_PORT")
	if serverPort == "" {
		log.Fatal("Err: HEALTHYAIR_SERVER_PORT is not set")
	}

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/healthyair?charset=utf8", healthyairSQLuser,
		healthyairSQLpassword))
	defer dbConn.Close()

	if err != nil {
		log.Fatal("Err on open database: ", err)
	}

	healthyairTARANTOOLserver = "127.0.0.1:3309"
	healthyairTARANTOOLopts = tarantool.Opts{
		Timeout:       50 * time.Millisecond,
		Reconnect:     100 * time.Millisecond,
		MaxReconnects: 3,
		User:          "guest",
	}
	healthyairTARANTOOLclient, err = tarantool.Connect(healthyairTARANTOOLserver, healthyairTARANTOOLopts)
	if err != nil {
		log.Fatal("@ERR ON CONNECTING TO TARANTOOL: ", err.Error())
	}

	resp, err := healthyairTARANTOOLclient.Ping()
	if err != nil {
		log.Println("Ping Code", resp.Code)
		log.Println("Ping Data", resp.Data)
		log.Fatal("@ERR ON PING: ", err)
	}

	// err = fasthttp.ListenAndServe("0.0.0.0:80", requestHandler)
	err = fasthttp.ListenAndServeTLS(":"+serverPort, certificatePath, keyPath,
		requestHandler)
	if err != nil {
		log.Fatal("@ERR ON STARTUP SERVER: ", err)
	}

}
