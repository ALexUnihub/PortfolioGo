package session

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	UserLogin string
}

//go:generate mockgen -source=session.go -destination=repo_mock.go -package=session SessionRepo
type SessionRepo interface {
	Create(w http.ResponseWriter, userID, userLogin string) (*Session, error)
	// Create(userID, userLogin string) (*Session, error)
	Check(r *http.Request) (*Session, error)
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func NewSession(userID, userLogin string) *Session {
	rand.Seed(time.Now().UnixNano())
	randID := string(randomBytes(10))

	return &Session{
		ID:        randID,
		UserID:    userID,
		UserLogin: userLogin,
	}
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomBytes(len int) []byte {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(97, 122))
	}
	return bytes
}
