package middleware

import (
	"bytes"
	"fmt"
	"myRedditClone/pkg/session"
	"net/http"
)

func Auth(sm *session.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		check := checkVote(r.URL.Path)

		if r.Method == "GET" && !check {
			next.ServeHTTP(w, r)
			return
		}

		// fmt.Println("auth middleware")
		sess, err := sm.Check(r)
		if err != nil {
			fmt.Println("no auth at", r.URL.Path, err.Error())
			http.Error(w, `no auth at`, http.StatusFound)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkVote(link string) bool {
	byteLink := []byte(link)
	idx := bytes.Index(byteLink, []byte("vote"))

	return idx != -1
}
