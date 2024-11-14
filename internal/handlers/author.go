package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)


type AuthorRepo interface{
	CreateAuthor(ctx context.Context, author *model.Author) error	
	GetAllAuthor()([]*model.Author, error)
	DeleteAuthorById(ctx context.Context, id int)error
	UpdateAuthor(ctx context.Context, id int,data *model.Author )error
} 

type AuthorHandler struct{
	Repo AuthorRepo
}

func NewAuthorHandler(mux *http.ServeMux, repo AuthorRepo){
	handler := &AuthorHandler{
		Repo:  repo,
	}

	mux.HandleFunc("POST /author", handler.CreateAuthor)
	mux.HandleFunc("PUT /author/{id}", handler.UpdateAuthor)
	mux.HandleFunc("DELETE /author/{id}", handler.DeleteAuthor)
	mux.HandleFunc("GET /author", handler.GetAllAuthor)
}

func (a *AuthorHandler) CreateAuthor(w http.ResponseWriter,r *http.Request){
	var author model.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Repo.CreateAuthor(r.Context(), &author); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author Created successfully"})
}

func (a *AuthorHandler) GetAllAuthor(w http.ResponseWriter, r *http.Request){
	authors,err := a.Repo.GetAllAuthor()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	jsonResp ,err := json.Marshal(authors)
	if err != nil{
		slog.Error("error marshal json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (a *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request){
	id,err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	if err := a.Repo.DeleteAuthorById(ctx,id); err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"message" : "Author successfully deleted"})
}

func (a *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request){
	id,err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	var author model.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Repo.UpdateAuthor(ctx,id,&author); err != nil {
		slog.Error("error using repo", "message",err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author updated successfully"})
}