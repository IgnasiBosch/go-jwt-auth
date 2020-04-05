package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	savePrivateKey("jwt_key.rsa", key)
	savePublicKey("jwt_key.pub", publicKey)
}

func savePrivateKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
