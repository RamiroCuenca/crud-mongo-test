package middlewares

import (
	"net/http"

	"github.com/RamiroCuenca/crud-mongo-test/auth"
	"github.com/RamiroCuenca/crud-mongo-test/common"
)

// Verifys that the user sent a valid JWT token
func Authenticated(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		_, err := auth.ValidateToken(token) // auth is the package we created
		// If token is invalid
		if err != nil {
			forbidden(w, r)
			return
		}

		f(w, r)
	}
}

// If the user is not authenticated it deprecates the request and returns a forbidden status
func forbidden(w http.ResponseWriter, r *http.Request) {
	json := []byte(`{
		"message": "It hasn't got authorization"
	}`)
	common.SendError(w, http.StatusForbidden, json)
}
