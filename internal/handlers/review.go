package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Fadil-Tao/manga-basis-data/internal/middleware"
	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type ReviewRepo interface{
	CreateReview(ctx context.Context, review *model.NewReviewRequest)error
	GetReviewFromManga(ctx context.Context, mangaId int)([]*model.Review, error)
	DeleteReview(ctx context.Context, doerId int, mangaId int, reviewerId int)error
	UpdateReview(ctx context.Context, doerId int,review *model.UpdateReview) error
	GetAReviewById(ctx context.Context, mangaId int, userId int)(*model.Review, error)
	ToggleLikeReview(ctx context.Context, doer int, mangaId int, reviewerId int)error
} 

type ReviewHandler struct {
	Repo ReviewRepo
}

func NewReviewHandler(mux *http.ServeMux, repo ReviewRepo){
	handler := & ReviewHandler{
		Repo: repo,
	}
	mux.Handle("POST /manga/{id}/review",middleware.Auth(http.HandlerFunc(handler.CreateReview)))
	mux.Handle("POST /manga/{mangaId}/review/{userId}/like",middleware.Auth(http.HandlerFunc(handler.ToggleLikeReview)))
	mux.Handle("DELETE /manga/{mangaId}/review/{userId}", middleware.Auth(http.HandlerFunc(handler.DeleteReview)))
	mux.Handle("PUT /manga/{mangaId}/review/{userId}", middleware.Auth(http.HandlerFunc(handler.UpdateReview)))
	mux.HandleFunc("GET /manga/{id}/review", handler.GetMangaReview)
	mux.HandleFunc("GET /manga/{mangaId}/review/{userId}", handler.GetReviewById)
}


func(rv *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	id :=r.PathValue("id")

	var review model.NewReviewRequest

	if err := json.NewDecoder(r.Body).Decode(&review);err != nil{
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	userIdstr := strconv.Itoa(userId)
	review.User_id = userIdstr
	review.Manga_id = id

	if err := rv.Repo.CreateReview(r.Context(), &review); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review Created successfully"})
}

func (rv *ReviewHandler) GetMangaReview(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	result, err := rv.Repo.GetReviewFromManga(ctx, id)
	if err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, statusCode(err))
		return
	}
	jsonResp, err := JSONMarshaller("review succesfully retrieved", result)
	if err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (rv *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request){
	mangaid , err := strconv.Atoi(r.PathValue("mangaId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	reviewerId , err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	if err := rv.Repo.DeleteReview(ctx,userId,mangaid,reviewerId); err !=nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Review successfully deleted"})
}
func (rv *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request){
	mangaid := r.PathValue("mangaId")
	reviewerId := r.PathValue("userId")

	var review model.UpdateReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	doerId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	review.User_id = reviewerId
	review.Manga_id = mangaid
	if err := rv.Repo.UpdateReview(ctx,doerId,&review); err != nil {
		slog.Error("error using repo", "message", err)
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return	
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review Updated successfully"})
}


func (rv *ReviewHandler) GetReviewById(w http.ResponseWriter, r *http.Request){
	mangaid , err := strconv.Atoi(r.PathValue("mangaId"))
	if err != nil {
		slog.Error("error converting mangaid to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	reviewerId , err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		slog.Error("error converting reviewerdi to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	review, err := rv.Repo.GetAReviewById(ctx,mangaid, reviewerId)
	if err != nil{
		slog.Error("error executing review query", "message", err)
		JSONError(w ,map[string]string{"message":err.Error()},statusCode(err))
		return
	}

	jsonResp, err := JSONMarshaller("review data retrieved succesfully", review)
	if err != nil{
		slog.Error("error marshalling", "message", err)
		JSONError(w ,map[string]string{"message":err.Error()},statusCode(err))
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (rv *ReviewHandler) ToggleLikeReview(w http.ResponseWriter, r *http.Request){
	mangaid , err := strconv.Atoi(r.PathValue("mangaId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	reviewerId , err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	doerId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	if err := rv.Repo.ToggleLikeReview(ctx, doerId,mangaid,reviewerId); err != nil {
		slog.Error("error executing review query", "message", err)
		JSONError(w ,map[string]string{"message":err.Error()},statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "review toggle liked triggered"})
}