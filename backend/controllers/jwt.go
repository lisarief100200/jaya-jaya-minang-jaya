package controllers

import (
	"backend/utils/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenGenerate is a function to generate token with data user. Return string and error.
func TokenAuthGenerator(username, level, sessionId, uid string) (string, time.Time, error) {
	// Tetapkan waktu kadaluarsa
	exp := time.Now().Add(24 * time.Hour)

	// Buat token dengan waktu kadaluarsa
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"uid":      uid,
		"exp":      exp.Unix(), // Waktu kadaluarsa dalam bentuk Unix timestamp
	})

	// Sign the token with a secret key.
	tokenString, err := token.SignedString([]byte(config.MyConfig.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, exp, nil
}
