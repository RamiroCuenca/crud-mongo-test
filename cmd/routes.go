package main

import (
	users "github.com/RamiroCuenca/vozy-test/users/controllers"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	// Set up the multiplexor
	mux := mux.NewRouter()

	// Set up routes
	mux.HandleFunc("/users/register", users.SignUp).Methods("POST")
	mux.HandleFunc("/users/login", users.SignIn).Methods("POST")
	mux.HandleFunc("/users/deletebyid", users.Delete).Methods("DELETE")

	return mux
}
