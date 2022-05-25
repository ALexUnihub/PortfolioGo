package session

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type SessionManager struct {
	DB          *sql.DB
	TokenSecret []byte
}

func NewSessionsManager(db *sql.DB) *SessionManager {
	var secretDB string
	tokenID := 1
	err := db.
		QueryRow("SELECT token FROM secret WHERE id = ?", tokenID).
		Scan(&secretDB)

	if err != nil {
		log.Println("pckg session, NewSessionsManager, scan err: ", err.Error())
		return nil
	}

	log.Println("check secret token: ", secretDB)

	return &SessionManager{
		DB:          db,
		TokenSecret: []byte(secretDB),
	}
}

func (sm *SessionManager) Create(w http.ResponseWriter, userID, userLogin string) (*Session, error) {
	// func (sm *SessionManager) Create(userID, userLogin string) (*Session, error) {
	session := &Session{}

	err := sm.DB.
		QueryRow("SELECT session_id FROM sessions WHERE (user_id = ?) AND (user_login = ?)", userID, userLogin).
		Scan(&session.ID)

	if err == nil {
		// found existed session in DB
		err = sendJWT(w, session.ID, userID, userLogin, sm.TokenSecret)
		if err != nil {
			log.Println("pckg session, Create, found existed session in DB, sendJWT", err.Error())
			return nil, err
		}

		session.UserID = userID
		session.UserLogin = userLogin
		log.Println("session found in DB for : sessionID", session.ID, "user ID", session.UserID, "user Login", session.UserLogin)
		return session, nil
	}

	session = NewSession(userID, userLogin)

	result, err := sm.DB.Exec(
		"INSERT INTO sessions (`session_id`, `user_id`, `user_login`) VALUES (?, ?, ?)",
		session.ID, session.UserID, session.UserLogin,
	)

	if err != nil {
		log.Println("pckg session, Create, sm.DB.Exec", err.Error())
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("pckg session, Create,  result.LastInsertId", err.Error())
		return nil, err
	}

	err = sendJWT(w, session.ID, session.UserID, session.UserLogin, sm.TokenSecret)
	if err != nil {
		log.Println("pckg session, Create, sendJWT", err.Error())
		return nil, err
	}

	log.Println("new session added in db: last ID", lastID)
	log.Println("new session created: session_id", session.ID, "user_id", session.UserID, "user_login", session.UserLogin)

	return session, nil
}

func (sm *SessionManager) Check(r *http.Request) (*Session, error) {
	inToken := r.Header.Get("authorization")
	inToken = inToken[len("Bearer "):]

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return sm.TokenSecret, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil {
		return nil, errors.New(`pckg session, Check, bad token`)
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(`pckg session, token.Claims, no payload`)
	}

	sess := &Session{}

	err = sm.DB.
		QueryRow("SELECT session_id, user_id, user_login FROM sessions WHERE session_id = ?", payload["sessionID"].(string)).
		Scan(&sess.ID, &sess.UserID, &sess.UserLogin)

	if err != nil {
		return nil, err
	}

	log.Println("session found for: sess_ID", sess.ID, "user ID:", sess.UserID, "user Login", sess.UserLogin)

	return sess, nil
}

func sendJWT(w http.ResponseWriter, sessionID, userID, userLogin string, secret []byte) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionID": sessionID,
		"user": map[string]interface{}{
			"username": userLogin,
			"id":       userID,
		},
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return errors.New(`pckg session, Create, token.SignedString err`)
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		return errors.New(`pckg session, Create, Marshal err`)
	}

	_, err = w.Write(resp)
	if err != nil {
		return errors.New(`pckg session, Create, Write err`)
	}

	return nil
}
