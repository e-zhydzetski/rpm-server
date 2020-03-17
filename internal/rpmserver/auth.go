package rpmserver

import (
	"net/http"
	"strings"
)

func NewHTTPAuthInterceptor(accessToken string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			auth = strings.TrimSpace(auth)
			authParts := strings.Split(auth, " ")
			if len(authParts) != 2 {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if authParts[0] != "Bearer" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			token := authParts[1]

			if token != accessToken {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
