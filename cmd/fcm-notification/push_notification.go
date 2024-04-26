package main

import (
	"encoding/base64"
	"os"
)

func getDecodeFireBaseKey() ([]byte, error) {
	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")

	decodeKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)

	if err != nil {
		return nil, err
	}
	return decodeKey, nil
}
