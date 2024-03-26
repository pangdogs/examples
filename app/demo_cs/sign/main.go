package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

func main() {
	priKey, err := rsa.GenerateKey(rand.Reader, sha256.Size*64)
	if err != nil {
		panic(err)
	}

	nowStr := time.Now().Format("2006-01-02T15_04_05")

	priKeyFile, err := os.Create(fmt.Sprintf("%s.pem", nowStr))
	if err != nil {
		panic(err)
	}
	defer priKeyFile.Close()

	err = pem.Encode(priKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priKey),
	})
	if err != nil {
		panic(err)
	}

	pubKeyFile, err := os.Create(fmt.Sprintf("%s.pub", nowStr))
	if err != nil {
		panic(err)
	}
	defer pubKeyFile.Close()

	err = pem.Encode(pubKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&priKey.PublicKey),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("saved to %s, %s\n", priKeyFile.Name(), pubKeyFile.Name())
}
