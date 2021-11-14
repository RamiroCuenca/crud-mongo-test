package main

import (
	"net/http"
	"sync"

	"github.com/RamiroCuenca/crud-mongo-test/routes"
)

var once sync.Once

func main() {
	// Get router
	mux := routes.GetRouter()

	// Run locally
	http.ListenAndServe(":8000", mux)
}
