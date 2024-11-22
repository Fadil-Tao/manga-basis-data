package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, newUser *model.NewUserRequest) error
	Login(ctx context.Context, user *model.UserLoginRequest) (*model.UserResponse, error)
}

type UserHandler struct {
	Repo UserRepository
}

func NewUserHandler(mux *http.ServeMux, repo UserRepository) {
	handler := &UserHandler{
		Repo: repo,
	}
	mux.HandleFunc("POST /auth/register", handler.RegisterUser)
	mux.HandleFunc("POST /auth/login", handler.Login)
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.NewUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}

	if err := u.Repo.RegisterUser(r.Context(), &user); err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Registered successfully"})
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	row, err := u.Repo.Login(r.Context(), &user)
	if err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	token, err := createToken(row)
	if err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name:     "token",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(168 * time.Hour),
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Message": "Login Success"})
}



func createToken(user *model.UserResponse) (string, error) {
	expiry := time.Now().Add(time.Hour * 168).Unix()
	secretKey := os.Getenv("JWT_SECRET")
	key := []byte(secretKey)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       user.Id,
			"is_admin": user.Is_admin,
			"exp":      expiry,
		})
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}
