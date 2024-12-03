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

type GenreRepo interface {
	CreateGenre(ctx context.Context, genre *model.Genre, userId int) error
	GetAllGenre() ([]*model.Genre, error)
	DeleteGenreById(ctx context.Context, id int, userId int) error
	UpdateGenre(ctx context.Context, id int, data *model.Genre, userId int) error
}

type GenreHandler struct {
	Repo GenreRepo
}

func NewGenreHandler(mux *http.ServeMux, repo GenreRepo) {
	handler := &GenreHandler{
		Repo: repo,
	}
	mux.Handle("POST /genre", middleware.Auth(http.HandlerFunc( handler.CreateGenre)))
	mux.Handle("PUT /genre/{id}",middleware.Auth( http.HandlerFunc(handler.UpdateGenre)))
	mux.Handle("DELETE /genre/{id}",middleware.Auth(http.HandlerFunc(handler.DeleteGenre)))
	mux.HandleFunc("GET /genre", handler.GetAllGenre)
}

func (g *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var genre model.Genre

	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
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

	if err := g.Repo.CreateGenre(r.Context(), &genre, userId); err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre Created successfully"})
}

func (g *GenreHandler) GetAllGenre(w http.ResponseWriter, r *http.Request) {
	genres, err := g.Repo.GetAllGenre()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResp, err := JSONMarshaller("Genres succesfully retrieved",genres)
	if err != nil {
		slog.Error("error marshal json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (g *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
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

	if err := g.Repo.DeleteGenreById(ctx, id, userId); err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre successfully deleted"})
}

func (g *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	var genre model.Genre
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}
	if err := g.Repo.UpdateGenre(ctx, id, &genre, userId); err != nil {
		slog.Error("error using repo", "message", err)
		JSONError(w, map[string]string{"message" :err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre updated successfully"})
}