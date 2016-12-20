package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
)

func main() {

	// Generate RSA Keys
	miryanPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	miryanPublicKey := &miryanPrivateKey.PublicKey

	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	raulPublicKey := &raulPrivateKey.PublicKey

	fmt.Println("Private Key : ", miryanPrivateKey)
	fmt.Println("Public key ", miryanPublicKey)
	fmt.Println("Private Key : ", raulPrivateKey)
	fmt.Println("Public key ", raulPublicKey)

}
