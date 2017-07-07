package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tarantool "github.com/tarantool/go-tarantool"
)

var healthyairSQLport, healthyairSQLuser, healthyairSQLpassword string
var dbConn *sql.DB

var healthyairTARANTOOLserver string
var healthyairTARANTOOLopts tarantool.Opts

var healthyairTARANTOOLclient *tarantool.Connection

func initial() {
	healthyairSQLuser = os.Getenv("HEALTHYAIR_SQL_USER")
	if healthyairSQLuser == "" {
		log.Fatal("Err: HEALTHYAIR_SQL_USER is not set")
	}
	healthyairSQLpassword = os.Getenv("HEALTHYAIR_SQL_PASSWORD")
	if healthyairSQLpassword == "" {
		log.Fatal("Err: HEALTHYAIR_SQL_PASSWORD is not set")
	}

	dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/healthyair?charset=utf8", healthyairSQLuser,
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

}

//authorize
func TestAuthorizeEnter(t *testing.T) {
	initial()

}

//register
func TestRegisterEnter(t *testing.T) {
	initial()
}

//sessions
func TestSessionStart(t *testing.T) {
	initial()
}
func TestSessionSetUserID(t *testing.T) {
	initial()
}
func TestSessionSetPreferredLanguage(t *testing.T) {
	initial()
}
func TestSessionResetEndTime(t *testing.T) {
	initial()
}
