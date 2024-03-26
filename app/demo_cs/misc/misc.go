package misc

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPublicKey(fp string) (*rsa.PublicKey, error) {
	bs, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bs)

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func LoadPrivateKey(fp string) (*rsa.PrivateKey, error) {
	bs, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bs)

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
