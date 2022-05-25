package user

type User struct {
	ID       string
	Login    string
	Password string
}

//go:generate mockgen -source=user.go -destination=repo_mock.go -package=user UserRepo
type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	AddNewUser(login, pass string) (*User, error)
}
