package middleware

import (
	"log"
	"net/http"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("access log middleware")
		next.ServeHTTP(w, r)
		log.Println("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
		)
	})
}
