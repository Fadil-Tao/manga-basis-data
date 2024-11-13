package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)


type AuthorRepo interface{
	CreateAuthor(ctx context.Context, author *model.Author) error	
} 

type AuthorHandler struct{
	Repo AuthorRepo
}

func NewAuthorHandler(mux *http.ServeMux, repo AuthorRepo){
	handler := &AuthorHandler{
		Repo:  repo,
	}

	mux.HandleFunc("POST /author", handler.CreateAuthor)
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