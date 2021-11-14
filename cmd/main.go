package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/RamiroCuenca/crud-mongo-test/auth"
	"github.com/RamiroCuenca/crud-mongo-test/routes"
)

var once sync.Once

func main() {
	dir, _ := os.Getwd()
	// Parse the certificates/keys (JWT)
	auth.LoadCertificates(
		dir+"/certificates/app.rsa",
		dir+"/certificates/app.rsa.pub",
	)
	// if err != nil {
	// 	log.Fatalf("Could not load the certificates/keys. Error: %v", err)
	// }

	// Get router
	mux := routes.GetRouter()

	// Run locally
	fmt.Println("Server running through port :8000 from container but exposed to :8080 to localhost...")
	http.ListenAndServe(":8000", mux)
}
