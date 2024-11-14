package middleware

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func AtleastAdmin(){
	
}