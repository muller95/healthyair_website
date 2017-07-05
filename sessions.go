package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	tarantool "github.com/tarantool/go-tarantool"
)

type session struct {
	sessionID         string
	userID            uint
	authorized        uint
	preferredLanguage string
	startTime         uint
}

func sessionStart() string {
	newSessID, err := uuid.NewRandom()
	if err != nil {
		log.Println("@ERR ON INITING SESSID")
		return ""
	}
	str := newSessID.String()
	_, err = healthyairTARANTOOLclient.Insert("healthyair", []interface{}{"aaaaa", 0, 0, "ru", uint(time.Now().Hour()*60 + time.Now().Minute())})
	if err != nil {
		log.Println("@ERR ON INSERT SESSION DATA TO TARANTOOL: ", err)
		return ""
	}
	return str
}
func sessionGet(SessID string) session {
	var uSession session
	resp, err := healthyairTARANTOOLclient.Select("healthyair", "sessionID", 0, 1, tarantool.IterEq, []interface{}{SessID})
	if err != nil {
		log.Println("Insert Error ", err)
		log.Println("Insert Code ", resp.Code)
		return uSession
	}

	uSession.sessionID = SessID
	uSession.userID = resp.Tuples()[0][1].(uint)
	uSession.authorized = resp.Tuples()[0][2].(uint)
	uSession.preferredLanguage = resp.Tuples()[0][3].(string)
	uSession.startTime = resp.Tuples()[0][4].(uint)

	return uSession
}
func sessionUpsert(Sess session) {
	_, err := healthyairTARANTOOLclient.Update("healthyair", "sessionID", []interface{}{Sess.sessionID}, []interface{}{
		[]interface{}{"=", 1, Sess.userID},
		[]interface{}{"=", 2, Sess.authorized},
		[]interface{}{"=", 3, Sess.preferredLanguage}})
	if err != nil {
		log.Println("@ERR ON UPSERT INFO TO TARANTOOL: ", err)
		return
	}
}
