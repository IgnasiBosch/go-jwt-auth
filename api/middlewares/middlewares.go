package middlewares

import (
	"errors"
	"github.com/IgnasiBosch/go-jwt-auth/api/auth"
	"github.com/IgnasiBosch/go-jwt-auth/api/responses"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareJWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.AccessTokenValid(r) {
			responses.JSONError(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		next(w, r)
	}
}
