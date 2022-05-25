package user

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"
)

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
	ErrNoNewUs = errors.New("user already exists")
)

type UserMemoryRepository struct {
	// lastRepoID int
	DB *sql.DB
	mu *sync.RWMutex
}

func NewMemoryRepo(db *sql.DB) *UserMemoryRepository {
	return &UserMemoryRepository{
		// lastRepoID: 2,
		DB: db,
		mu: &sync.RWMutex{},
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	user := &User{}
	repo.mu.RLock()

	err := repo.DB.
		QueryRow("SELECT user_id, login, password FROM users WHERE login = ?", login).
		Scan(&user.ID, &user.Login, &user.Password)

	repo.mu.RUnlock()

	if err != nil {
		return nil, ErrNoUser
	}

	if user.Password != pass {
		return nil, ErrBadPass
	}

	return user, nil
}

func (repo *UserMemoryRepository) AddNewUser(login, pass string) (*User, error) {
	var checkLogin string
	var usLastID int

	repo.mu.RLock()

	err := repo.DB.
		QueryRow("SELECT login FROM users WHERE login = ?", login).
		Scan(&checkLogin)

	repo.mu.RUnlock()

	if err == nil {
		log.Println(ErrNoNewUs.Error(), "checkLogin", checkLogin, "login", login)
		return nil, ErrNoNewUs
	}

	err = repo.DB.
		QueryRow("SELECT MAX(id) AS max_id FROM users").
		Scan(&usLastID)

	if err != nil {
		log.Println("pckg user, AddNewUser, search max id err:", err.Error())
		return nil, err
	}

	repo.mu.Lock()

	usID := strconv.Itoa(usLastID + 1)

	result, err := repo.DB.Exec(
		"INSERT INTO users (`user_id`, `login`, `password`) VALUES (?, ?, ?)",
		usID, login, pass,
	)

	repo.mu.Unlock()

	if err != nil {
		log.Println("pckg user, AddNewUser, repo.DB.Exec, err : ", err.Error())
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("pckg user, AddNewUser, result.LastInsertId(), err : ", err.Error())
		return nil, err
	}
	log.Println("last id in MySQL: ", lastID)

	newUser := &User{
		ID:       usID,
		Login:    login,
		Password: pass,
	}

	log.Println("user added: ID", newUser.ID, "login", newUser.Login)
	return newUser, nil
}
