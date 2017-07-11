package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	tarantool "github.com/tarantool/go-tarantool"
)

type Session struct {
	SessionID         string
	UserID            uint64
	Authorized        bool
	PreferredLanguage string
	EndTime           int64
}

var sessionsSpace = "sessions"
var sessionLifeTime int64 = 3 * 24 * 60 * 60

func SessionStart(language string) (Session, RestCode) {
	var session Session

	newsid, err := uuid.NewRandom()
	if err != nil {
		log.Println("Err on creating session id: err")
		return session, InternalServerError
	}
	str := newsid.String()
	session.SessionID = str
	session.UserID = 0
	session.Authorized = false
	session.PreferredLanguage = language
	session.EndTime = time.Now().Unix() + sessionLifeTime

	resp, err := healthyairTARANTOOLclient.Insert(sessionsSpace, []interface{}{str, 0, false, language,
		time.Now().Unix() + sessionLifeTime})
	if err != nil {
		log.Printf("Err insert session: %v %v\n", err, resp.Code)
		return session, InternalServerError
	}
	return session, Ok
}

func SessionSetUserID(session *Session) {
	_, err := healthyairTARANTOOLclient.Update(sessionsSpace, "primary", []interface{}{session.SessionID},
		[]interface{}{[]interface{}{"=", 1, session.UserID}, []interface{}{"=", 2, session.Authorized}})
	if err != nil {
		log.Println("Err on upsert userID into tarantool: ", err)
		return
	}
}
func SessionSetPreferredLanguage(session *Session) RestCode {
	_, err := healthyairTARANTOOLclient.Update(sessionsSpace, "primary", []interface{}{session.SessionID},
		[]interface{}{[]interface{}{"=", 3, session.PreferredLanguage}})
	if err != nil {
		log.Println("Err on upsert language into tarantool: ", err)
		return InternalServerError
	}

	return Ok
}
func SessionResetEndTime(session *Session) {
	_, err := healthyairTARANTOOLclient.Update(sessionsSpace, "primary", []interface{}{session.SessionID},
		[]interface{}{[]interface{}{"=", 4, session.EndTime}})
	if err != nil {
		log.Println("Err on upsert session into tarantool: ", err)
		return
	}
}

func SessionDelete(sid string) {
	resp, err := healthyairTARANTOOLclient.Delete(sessionsSpace, "primary", []interface{}{sid})
	if err != nil {
		log.Printf("Err delete session: %v %v\n", err, resp.Code)
		return
	}
}

func SessionGet(sid string) (Session, RestCode) {
	var session Session
	resp, err := healthyairTARANTOOLclient.Select(sessionsSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{sid})
	if err != nil {
		log.Printf("Err select session: %v %v\n", err, resp.Code)
		return session, InternalServerError
	}

	if len(resp.Tuples()) == 0 {
		return session, NotFound
	}

	session.SessionID = sid
	session.UserID = resp.Tuples()[0][1].(uint64)
	session.Authorized = resp.Tuples()[0][2].(bool)
	session.PreferredLanguage = resp.Tuples()[0][3].(string)
	session.EndTime = int64(resp.Tuples()[0][4].(uint64))

	log.Println(session)
	if time.Now().Unix() > session.EndTime {
		SessionDelete(sid)
		return session, SessionExpired
	}

	return session, Ok
}

func SessionUpsert(Sess *Session) {
	_, err := healthyairTARANTOOLclient.Update(sessionsSpace, "primary", []interface{}{Sess.SessionID},
		[]interface{}{[]interface{}{"=", 1, Sess.UserID}, []interface{}{"=", 2, Sess.Authorized},
			[]interface{}{"=", 3, Sess.PreferredLanguage}})
	if err != nil {
		log.Println("Err on upsert session into tarantool: ", err)
		return
	}
}
