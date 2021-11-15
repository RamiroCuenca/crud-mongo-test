package routes

import (
	"github.com/RamiroCuenca/crud-mongo-test/middlewares"
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
	mux.HandleFunc("/users/deletebyid", middlewares.Authenticated(users.Delete)).Methods("DELETE")
	mux.HandleFunc("/users/updatebyid", middlewares.Authenticated(users.Update)).Methods("PUT")

	// mux.HandleFunc("/users/deletebyid", users.Delete).Methods("DELETE")
	// mux.HandleFunc("/users/updatebyid", users.Update).Methods("PUT")

	// Set up routes for posts
	mux.HandleFunc("/posts/create", middlewares.Authenticated(posts.Create)).Methods("POST")
	mux.HandleFunc("/posts/getbyid", middlewares.Authenticated(posts.GetById)).Methods("GET")
	mux.HandleFunc("/posts/getallfromuserid", middlewares.Authenticated(posts.GetAllFromUserId)).Methods("GET")
	mux.HandleFunc("/posts/deletebyid", middlewares.Authenticated(posts.Delete)).Methods("DELETE")
	mux.HandleFunc("/posts/updatebyid", middlewares.Authenticated(posts.Update)).Methods("PUT")

	return mux
}
