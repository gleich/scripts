package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mattglei.ch/timber"
)

func generateAppToken(inputs Inputs) (string, time.Time) {
	key, err := os.ReadFile(inputs.KeyFilename)
	if err != nil {
		timber.Fatal(err, "failed to read from", inputs.KeyFilename)
	}

	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PRIVATE KEY" {
		timber.FatalMsg("failed to decode PEM block containing private key")
	}

	keyIfc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		timber.Fatal(err, "parse PKCS#8 private key", err)
	}
	ecKey, ok := keyIfc.(*ecdsa.PrivateKey)
	if !ok {
		timber.Fatal(err, "not an ECDSA private key")
	}

	now := time.Now()
	expiration := now.Add(182 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"iss": inputs.TeamID,
		"iat": now.Unix(),
		"exp": expiration.Unix(), // expires in 6 months
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = inputs.KeyID
	token.Header["alg"] = "ES256"

	appToken, err := token.SignedString(ecKey)
	if err != nil {
		timber.Fatal(err, "failed to sign token")
	}

	return appToken, expiration
}
