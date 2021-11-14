package main

import (
	users "github.com/RamiroCuenca/vozy-test/users/controllers"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	// Set up the multiplexor
	mux := mux.NewRouter()

	// Set up routes
	mux.HandleFunc("/users/create", users.Create).Methods("POST")

	return mux
}
