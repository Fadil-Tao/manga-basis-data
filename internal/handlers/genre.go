package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type GenreRepo interface{
	CreateGenre(ctx context.Context,genre *model.Genre)error
	GetAllGenre()([]*model.Genre, error)
	DeleteGenreById(ctx context.Context, id int)error
	UpdateGenre(ctx context.Context, id int,data *model.Genre )error
}

type GenreHandler struct {
	Repo GenreRepo
}

func NewGenreHandler(mux *http.ServeMux, repo GenreRepo){
	handler := &GenreHandler{
		Repo:  repo,
	}
	mux.HandleFunc("POST /genre", handler.CreateGenre)
	mux.HandleFunc("PUT /genre/{id}", handler.UpdateGenre)
	mux.HandleFunc("DELETE /genre/{id}", handler.DeleteGenre)
	mux.HandleFunc("GET /genre", handler.GetAllGenre)
}

func (g *GenreHandler) CreateGenre(w http.ResponseWriter , r *http.Request){
	var genre model.Genre
	
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := g.Repo.CreateGenre(r.Context(), &genre); err != nil {
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre Created successfully"})
}

func (g *GenreHandler) GetAllGenre(w http.ResponseWriter, r *http.Request){
	genres,err := g.Repo.GetAllGenre()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	jsonResp ,err := json.Marshal(genres)
	if err != nil{
		slog.Error("error marshal json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (g *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request){
	id,err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	if err := g.Repo.DeleteGenreById(ctx,id); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"message" : "Genre successfully deleted"})
}

func (g *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request){
	id,err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	var genre model.Genre
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := g.Repo.UpdateGenre(ctx,id,&genre); err != nil {
		slog.Error("error using repo", "message",err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre updated successfully"})
}
