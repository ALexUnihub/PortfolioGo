package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"myRedditClone/pkg/session"
	"myRedditClone/pkg/user"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandlerIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	service := &UserHandler{
		UserRepo: st,
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	service.Index(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	checkHTML := `doctype html`
	if !bytes.Contains(body, []byte(checkHTML)) {
		t.Errorf("no HTML found")
		return
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctrlSess := gomock.NewController(t)
	defer ctrlSess.Finish()

	sm := session.NewMockSessionRepo(ctrlSess)
	st := user.NewMockUserRepo(ctrl)
	usHand := &UserHandler{
		UserRepo: st,
		Sessions: sm,
	}
	resultUser := []*user.User{
		{ID: "us1", Login: "12", Password: "11112222"},
	}

	resultSess := []*session.Session{
		{ID: "1", UserID: "us1", UserLogin: "12"},
	}

	// correct input
	bodyReader := strings.NewReader(`{"username": "12", "password": "11112222"}`)

	req := httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	st.EXPECT().AddNewUser(resultUser[0].Login, resultUser[0].Password).Return(resultUser[0], nil)
	sm.EXPECT().Create(w, resultUser[0].ID, resultUser[0].Login).Return(resultSess[0], nil)

	usHand.Register(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(``)) {
		t.Errorf("unexpected error")
		return
	}

	// unknown payload (no set header)
	req = httptest.NewRequest("POST", "/api/register", nil)
	w = httptest.NewRecorder()

	usHand.Register(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`unknown payload`)) {
		t.Errorf("unexpected error: expect body : unknown payload")
		return
	}

	// bad JSON
	bodyReader = strings.NewReader(`}`)
	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	usHand.Register(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`cant unpack payload`)) {
		t.Errorf("unexpected error: expect body : unknown payload")
		return
	}

	// add new user - error
	bodyReader = strings.NewReader(`{"username": "12", "password": "11112222"}`)
	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	st.EXPECT().AddNewUser(resultUser[0].Login, resultUser[0].Password).Return(nil, fmt.Errorf("add new user err"))

	usHand.Register(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`add new user err`)) {
		t.Errorf("unexpected error: expect body : add new user err")
		return
	}

	// create session - error
	bodyReader = strings.NewReader(`{"username": "12", "password": "11112222"}`)
	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	st.EXPECT().AddNewUser(resultUser[0].Login, resultUser[0].Password).Return(resultUser[0], nil)
	sm.EXPECT().Create(w, resultUser[0].ID, resultUser[0].Login).Return(nil, fmt.Errorf("create session err"))

	usHand.Register(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`create session err`)) {
		t.Errorf("unexpected error: expect body : create session err")
		return
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctrlSess := gomock.NewController(t)
	defer ctrlSess.Finish()

	sm := session.NewMockSessionRepo(ctrlSess)
	st := user.NewMockUserRepo(ctrl)
	usHand := &UserHandler{
		UserRepo: st,
		Sessions: sm,
	}
	resultUser := []*user.User{
		{ID: "us1", Login: "12", Password: "11112222"},
	}

	resultSess := []*session.Session{
		{ID: "1", UserID: "us1", UserLogin: "12"},
	}

	// correct input
	bodyReader := strings.NewReader(`{"username": "12", "password": "11112222"}`)

	req := httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	st.EXPECT().Authorize(resultUser[0].Login, resultUser[0].Password).Return(resultUser[0], nil)
	sm.EXPECT().Create(w, resultUser[0].ID, resultUser[0].Login).Return(resultSess[0], nil)

	usHand.Login(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(``)) {
		t.Errorf("unexpected error")
		return
	}

	// unknown payload
	req = httptest.NewRequest("POST", "/api/register", nil)
	w = httptest.NewRecorder()

	usHand.Login(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`unknown payload`)) {
		t.Errorf("unexpected error, expected : unknown payload")
		return
	}

	// cant unpack payload
	bodyReader = strings.NewReader(`}`)

	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	usHand.Login(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`cant unpack payload`)) {
		t.Errorf("unexpected error, expected : cant unpack payload")
		return
	}

	// Authorize - err
	bodyReader = strings.NewReader(`{"username": "12", "password": "11112222"}`)

	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	st.EXPECT().Authorize(resultUser[0].Login, resultUser[0].Password).Return(nil, fmt.Errorf("Authorize err"))

	usHand.Login(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`Authorize err`)) {
		t.Errorf("unexpected error, expected : Authorize err")
		return
	}

	// Create session - err
	bodyReader = strings.NewReader(`{"username": "12", "password": "11112222"}`)

	req = httptest.NewRequest("POST", "/api/register", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	st.EXPECT().Authorize(resultUser[0].Login, resultUser[0].Password).Return(resultUser[0], nil)
	sm.EXPECT().Create(w, resultUser[0].ID, resultUser[0].Login).Return(nil, fmt.Errorf("Create session - err"))

	usHand.Login(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`Create session - err`)) {
		t.Errorf("unexpected error, expected : Create session - err")
		return
	}
}
