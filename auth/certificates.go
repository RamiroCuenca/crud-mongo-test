package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)

// We left them private (Only can be used by this package)
// so that the other packages can not manipulate them

// Singleton that assign a value to this vars only once (EXPORTED)
func LoadCertificates() error {
	var err error

	once.Do(func() {
		err = loadCertificates()
	})

	return err
}

// Loads the certicates and send them to be parsed
func loadCertificates() error {
	// dir, _ := os.Getwd()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}

	privateBytes, err := ioutil.ReadFile("./app.rsa")
	if err != nil {
		return err
	}

	publicBytes, err := ioutil.ReadFile("./app.rsa.pub")
	if err != nil {
		return err
	}

	return parseRSA(privateBytes, publicBytes)
}

// Parces the certicates
func parseRSA(privateBytes, publicBytes []byte) error {
	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
