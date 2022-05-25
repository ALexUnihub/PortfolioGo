package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"myRedditClone/pkg/session"
	"myRedditClone/pkg/user"
)

type UserHandler struct {
	UserRepo user.UserRepo
	// Sessions *session.SessionManager
	Sessions session.SessionRepo
}

type LoginForm struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type MessageUser struct {
	Msg string `json:"message"`
}

// показ начальной страницы
func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	dataHTML, err := os.Open("../../static/html/index.html")
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}

	byteValue, err := ioutil.ReadAll(dataHTML)
	if err != nil {
		http.Error(w, `Parsing html error`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteValue)
	if err != nil {
		log.Println("err in pckg handlers, Index", err.Error())
		return
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, `unknown payload`, http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	fd := &LoginForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		http.Error(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	us, err := h.UserRepo.AddNewUser(fd.Login, fd.Password)
	if err != nil {
		msg := &MessageUser{
			Msg: err.Error(),
		}
		w.WriteHeader(401)
		err = sendJSONUser(w, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
	// create session in db / JWT
	sess, err := h.Sessions.Create(w, us.ID, us.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Register", "created session for:", us.Login, "user ID:", sess.UserID, "session ID:", sess.ID)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, `unknown payload`, http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	fd := &LoginForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		http.Error(w, `cant unpack payload`, http.StatusBadRequest)
		return
	}

	us, err := h.UserRepo.Authorize(fd.Login, fd.Password)
	if err != nil {
		msg := &MessageUser{
			Msg: err.Error(),
		}
		w.WriteHeader(401)
		err = sendJSONUser(w, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
	// create session login / JWT
	sess, err := h.Sessions.Create(w, us.ID, us.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Login", "created session for:", us.Login, "user ID:", sess.UserID, "session ID:", sess.ID)
}

func sendJSONUser(w http.ResponseWriter, data interface{}) error {
	byteValue, err := json.Marshal(data)
	if err != nil {
		return errors.New("sendJSON, json.Marshal err")
	}

	_, err = w.Write(byteValue)
	if err != nil {
		return errors.New("sendJSON, w.Write err")
	}

	return nil
}
