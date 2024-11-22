package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type key string

const userKey key = "user"

func VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		tokenValue := cookie.Value
		token, err := VerifyToken(tokenValue)
		if err != nil {
			http.Error(w, "Unauthorized invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := r.Context()

			ctx = context.WithValue(ctx, userKey, claims)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AtleastAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userKey).(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid user claims", http.StatusUnauthorized)
			return
		}

		isAdmin, ok := claims["is_admin"]
		if !ok || isAdmin != 1 {
			http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	claims, ok := r.Context().Value(userKey).(jwt.MapClaims)
	if !(ok) {
		return 0, fmt.Errorf("unauthorized: invalid user claims")
	}

	userIdFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unauthorized: missing or invalid user ID")
	}
	userId := int(userIdFloat)
	return userId, nil
}
