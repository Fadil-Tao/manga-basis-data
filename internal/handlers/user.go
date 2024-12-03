package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Fadil-Tao/manga-basis-data/internal/middleware"
	"github.com/Fadil-Tao/manga-basis-data/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, newUser *model.NewUserRequest) error
	Login(ctx context.Context, user *model.UserLoginRequest) (*model.UserResponse, error)
	SearchUser(ctx context.Context, username string)([]*model.UserResponse, error)
	UpdateUser(ctx context.Context, usernameTarget string,userid int,userData *model.NewUserRequest) error
	GetUserDetailByUsername(ctx context.Context, username string)(*model.UserResponse, error)
	GetUserRatedManga(ctx context.Context, username string)([]*model.UserRatedManga, error)
	GetUserLikedManga(ctx context.Context, username string)([]*model.UserLikedManga, error)
	GetUserReadlist(ctx context.Context, username string) ([]*model.Readlist, error)
	DeleteUser(ctx context.Context, userId int,username string)error
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
	mux.HandleFunc("DELETE /auth/logout", handler.Logout)
	mux.Handle("PUT /users/{username}", middleware.Auth(http.HandlerFunc(handler.UpdateUser)))
	mux.Handle("DELETE /users/{username}", middleware.Auth(http.HandlerFunc(handler.DeleteUser)))
	mux.HandleFunc("GET /users",handler.SearchUser)
	mux.HandleFunc("GET /users/{username}", handler.GetUserDetailByUsername)
	mux.HandleFunc("GET /users/{username}/likedmanga", handler.GetUserLikedManga)
	mux.HandleFunc("GET /users/{username}/ratedmanga", handler.GetUserRateddManga)
	mux.HandleFunc("GET /users/{username}/readlist", handler.GetUserReadlist)	
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.NewUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}

	if err := u.Repo.RegisterUser(r.Context(), &user); err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, statusCode(err))
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
		Path:     "/",
		Expires:  time.Now().Add(168 * time.Hour),
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Message": "Login Success"})
}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request){
	c := http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Message": "Logout Success"})
}

func createToken(user *model.UserResponse) (string, error) {
	expiry := time.Now().Add(time.Hour * 168).Unix()
	secretKey := os.Getenv("JWT_SECRET")
	key := []byte(secretKey)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       user.Id,
			"exp":      expiry,
		})
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter , r *http.Request){
	username := r.PathValue("username")	

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	var user model.NewUserRequest 
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil{
		JSONError(w, map[string]string{"message": "bad request"}, http.StatusBadRequest)
		return
	}

	if err := u.Repo.UpdateUser(r.Context(), username,userId,&user); err != nil{
		JSONError(w, map[string]string{"message" : err.Error()}, statusCode(err)) 
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Updated successfully"})	
}


func (u *UserHandler) GetUserDetailByUsername(w http.ResponseWriter, r *http.Request){
	username := r.PathValue("username")
	
	userdetail,err := u.Repo.GetUserDetailByUsername(r.Context(), username) ;
	if err != nil{
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return	
	}


	jsonResp, err := JSONMarshaller("user succesfully retrieved", userdetail)
	if err != nil {
		JSONError(w, map[string]string{"messaage": "internal server error"}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)	
} 

func (u *UserHandler) SearchUser(w http.ResponseWriter, r *http.Request){
	var username string
	
	if query := r.URL.Query().Get("username"); query != "" {
		username = query
	}
	users , err := u.Repo.SearchUser(r.Context(),username)
	if err != nil {
		JSONError(w, map[string]string{"message" :err.Error()},statusCode(err) )
		return
	}

	jsonResp, err := JSONMarshaller("user succesfully retrieved", users)
	if err != nil {
		JSONError(w, map[string]string{"messaage": "internal server error"}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)	
}


func (u *UserHandler) GetUserLikedManga(w http.ResponseWriter, r *http.Request){
	username := r.PathValue("username")
	
	mangas, err := u.Repo.GetUserLikedManga(r.Context(), username)
	if err != nil {
		JSONError(w, map[string]string{"message" :err.Error()},statusCode(err) )
		return
	}
	jsonResp, err := JSONMarshaller("mangas succesfully retrieved" , mangas)
	if err != nil {
		JSONError(w, map[string]string{"message" :"internal server error"},http.StatusInternalServerError )
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)	
}

func (u *UserHandler) GetUserRateddManga(w http.ResponseWriter, r *http.Request){
	username := r.PathValue("username")
	
	mangas, err := u.Repo.GetUserRatedManga(r.Context(), username)
	if err != nil {
		JSONError(w, map[string]string{"message" :err.Error()},statusCode(err) )
		return
	}
	jsonResp, err := JSONMarshaller("mangas succesfully retrieved" , mangas)
	if err != nil {
		JSONError(w, map[string]string{"message" :"internal server error"},http.StatusInternalServerError )
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)	
}

func (u *UserHandler) GetUserReadlist(w http.ResponseWriter, r *http.Request){
	username := r.PathValue("username")
	
	readlist, err := u.Repo.GetUserReadlist(r.Context(), username)
	if err != nil {
		JSONError(w, map[string]string{"message" :err.Error()},statusCode(err) )
		return
	}

	jsonResp, err := JSONMarshaller("mangas succesfully retrieved" , readlist)
	if err != nil {
		JSONError(w, map[string]string{"message" :"internal server error"},http.StatusInternalServerError )
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)	
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	username := r.PathValue("username")
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	
	if err = u.Repo.DeleteUser(r.Context(),userId, username); err != nil{
		JSONError(w, map[string]string{
			"message": err.Error(),
		},statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Message": "user deleted Successfully"})
}

