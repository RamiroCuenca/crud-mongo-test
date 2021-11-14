package routes

import (
	posts "github.com/RamiroCuenca/crud-mongo-test/posts/controllers"
	users "github.com/RamiroCuenca/crud-mongo-test/users/controllers"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	// Set up the multiplexor
	mux := mux.NewRouter()

	// Set up routes for users controllers
	mux.HandleFunc("/users/register", users.SignUp).Methods("POST")
	mux.HandleFunc("/users/login", users.SignIn).Methods("POST")
	// mux.HandleFunc("/users/deletebyid", middlewares.Authenticated(users.Delete)).Methods("DELETE")
	// mux.HandleFunc("/users/updatebyid", middlewares.Authenticated(users.Update)).Methods("PUT")

	mux.HandleFunc("/users/deletebyid", users.Delete).Methods("DELETE")
	mux.HandleFunc("/users/updatebyid", users.Update).Methods("PUT")

	// Set up routes for posts
	mux.HandleFunc("/posts/create", posts.Create).Methods("POST")
	mux.HandleFunc("/posts/getbyid", posts.GetById).Methods("GET")

	return mux
}
