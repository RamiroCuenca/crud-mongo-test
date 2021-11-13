package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/RamiroCuenca/vozy-test/database"
)

var once sync.Once

func main() {
	client := database.GetMongoClient()

	once.Do(func() {
		if client != nil {
			fmt.Println("Succesfully connected to database!")
		}
	})

	// Run locally
	http.ListenAndServe(":8000", nil)
}
