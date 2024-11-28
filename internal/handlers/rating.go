package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Fadil-Tao/manga-basis-data/internal/middleware"
)


type RatingRepo interface{
	RateManga(ctx context.Context, userId int, mangaId int, rating int)error
}

type RatingHandler struct {
	Repo RatingRepo
}

func NewRatingHandler(mux *http.ServeMux, repo RatingRepo){
	handlers := &RatingHandler{
		Repo: repo,
	}

	mux.Handle("POST /rating",middleware.Auth(http.HandlerFunc(handlers.RateManga)))
}


func (rm *RatingHandler) RateManga(w http.ResponseWriter, r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	ratingRequest := struct{
		MangaId string `json:"mangaId"`
		Rating int `json:"rating"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&ratingRequest); err != nil{
		JSONError(w,map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}
	mangaId, err := strconv.Atoi(ratingRequest.MangaId)
	if err != nil {
		JSONError(w, map[string]string{"message": "internal server error"}, http.StatusInternalServerError)
		return 
	}
	slog.Info("manga id request",":", mangaId)

	if err := rm.Repo.RateManga(r.Context(),userId,mangaId,ratingRequest.Rating); err != nil{
		JSONError(w,map[string]string{
			"message": err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Succesfully rate manga"})
}


