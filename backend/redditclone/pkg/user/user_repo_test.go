package user

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.NewRows([]string{"user_id", "login", "password"})
	expect := []*User{
		{"userID", "userLogin", "userPass"},
		{"userID", "userLogin", "userPass"},
		{"userID", "userLogin", "userPass"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Login, item.Password)
	}

	repo := &UserMemoryRepository{
		DB: db,
		mu: &sync.RWMutex{},
	}

	// return user
	mock.
		ExpectQuery("SELECT user_id, login, password FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	user, err := repo.Authorize(expect[0].Login, expect[0].Password)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(user, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], user)
		return
	}

	// return ErrNoUser
	mock.
		ExpectQuery("SELECT user_id, login, password FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnError(ErrNoUser)

	_, err = repo.Authorize(expect[0].Login, expect[0].Password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// return ErrNoUser
	mock.
		ExpectQuery("SELECT user_id, login, password FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	_, err = repo.Authorize(expect[0].Login, "wrong_pass")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// return ErrNoUser
	mock.
		ExpectQuery("SELECT user_id, login, password FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	_, err = repo.Authorize(expect[0].Login, "wrong_pass")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}

func TestAddNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &UserMemoryRepository{
		DB: db,
		mu: &sync.RWMutex{},
	}

	// good query
	rows := sqlmock.NewRows([]string{"login"}).
		AddRow("userLogin")
	expect := []*User{
		{"3", "userLogin", "userPass"},
		{"4", "userLogin", "userPass"},
	}

	// ErrNoNewUs - user already exists
	mock.
		ExpectQuery("SELECT login FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	_, err = repo.AddNewUser(expect[0].Login, expect[0].Password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// add new user
	rowID := sqlmock.NewRows([]string{"max_id"})
	rowID = rowID.AddRow(2)
	// expect[0].ID = "3"

	mock.
		ExpectQuery("SELECT login FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnError(fmt.Errorf("no user in DB"))

	mock.
		ExpectQuery(regexp.QuoteMeta("SELECT MAX(id) AS max_id FROM users")).
		WithArgs().
		WillReturnRows(rowID)

	mock.
		ExpectExec("INSERT INTO users").
		WithArgs(expect[0].ID, expect[0].Login, expect[0].Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := repo.AddNewUser(expect[0].Login, expect[0].Password)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(user, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], user)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// select max_id err
	mock.
		ExpectQuery("SELECT login FROM users WHERE").
		WithArgs(expect[0].Login).
		WillReturnError(fmt.Errorf("no user in DB"))

	mock.
		ExpectQuery(regexp.QuoteMeta("SELECT MAX(id) AS max_id FROM users")).
		WithArgs().
		WillReturnError(fmt.Errorf("max id search err"))

	_, err = repo.AddNewUser(expect[0].Login, expect[0].Password)
	if err == nil {
		t.Errorf("unexpected err: got nil, expected: max id search err")
		return
	}
	if err.Error() != "max id search err" {
		t.Errorf("unexpected err, expected: max id search err, got: %s", err.Error())
		return
	}

	// repo.DB.Exec, err
	rowID = rowID.AddRow(3)
	mock.
		ExpectQuery("SELECT login FROM users WHERE").
		WithArgs(expect[1].Login).
		WillReturnError(fmt.Errorf("bd - err"))

	mock.
		ExpectQuery(regexp.QuoteMeta("SELECT MAX(id) AS max_id FROM users")).
		WithArgs().
		WillReturnRows(rowID)

	mock.
		ExpectExec("INSERT INTO users").
		WithArgs(expect[1].ID, expect[1].Login, expect[1].Password).
		WillReturnError(fmt.Errorf("DB err"))

	_, err = repo.AddNewUser(expect[1].Login, expect[1].Password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}
