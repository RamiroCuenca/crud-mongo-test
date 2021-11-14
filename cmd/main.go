package main

import (
	"net/http"
	"sync"
)

var once sync.Once

func main() {
	// Get router
	mux := GetRouter()

	// Run locally
	http.ListenAndServe(":8000", mux)
}
