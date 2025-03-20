package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("JWT123")

func GenerateStudentJWT(userID int, userType string, blockID uint) (string, error) {
	// Create a new JWT token with user payload
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = map[string]interface{}{
		"user_id":  userID,
		"type":     userType,
		"block_id": blockID,
	}
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	// Sign the token with the secret key and return
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func DecodeStudentJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract and return the claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
