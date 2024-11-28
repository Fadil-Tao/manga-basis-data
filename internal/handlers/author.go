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

type AuthorRepo interface {
	CreateAuthor(ctx context.Context, author *model.Author, userId int) error
	GetAllAuthor() ([]*model.Author, error)
	DeleteAuthorById(ctx context.Context, id int, userId int) error
	UpdateAuthor(ctx context.Context, id int, data *model.Author, userId int) error
	SearchAuthor(ctx context.Context, input string) ([]model.Author,error)
	GetAuthorById(ctx context.Context, id int) (*model.Author, error) 
	GetAuthorManga(ctx context.Context, id int)([]model.Manga, error)
}

type AuthorHandler struct {
	Repo AuthorRepo
}

func NewAuthorHandler(mux *http.ServeMux, repo AuthorRepo) {
	handler := &AuthorHandler{
		Repo: repo,
	}

	mux.Handle("POST /author", middleware.Auth(http.HandlerFunc(handler.CreateAuthor)))
	mux.Handle("PUT /author/{id}", middleware.Auth(http.HandlerFunc(handler.UpdateAuthor)))
	mux.Handle("DELETE /author/{id}", middleware.Auth(http.HandlerFunc(handler.DeleteAuthor)))
	mux.HandleFunc("GET /author", handler.GetAllAuthor)
	mux.HandleFunc("GET /author/{id}/manga", handler.GetAllAuthor)
	mux.HandleFunc("GET /author/{id}", handler.GetAuthorDetails)
}

func (a *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author model.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	if err := a.Repo.CreateAuthor(r.Context(), &author, userId); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author Created successfully"})
}

func (a *AuthorHandler) GetAllAuthor(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	if query != "" {
		ctx := r.Context()
		authors, err := a.Repo.SearchAuthor(ctx,query) 
		if err != nil{
			slog.Error("error at calling procedure", "message", err)
			JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
			return	
		}
		jsonResp, err := json.Marshal(authors)
		if err != nil {
			slog.Error("error marshal json")
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return
	}
	authors, err := a.Repo.GetAllAuthor()
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	jsonResp, err := JSONMarshaller("Authors Succesfully Retrieved",authors)
	if err != nil {
		slog.Error("error marshal json")
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (a *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
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

	if err := a.Repo.DeleteAuthorById(ctx, id, userId); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Author successfully deleted"})
}

func (a *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting to int json")
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	var author model.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := a.Repo.UpdateAuthor(ctx, id, &author, userId); err != nil {
		slog.Error("error using repo", "message", err)
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author updated successfully"})
}

func (a *AuthorHandler) GetAuthorDetails(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil{
		slog.Error("error converting to int json")
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	author, err := a.Repo.GetAuthorById(ctx,id)
	if err != nil{
		slog.Error("error executing author query", "message", err)
		JSONError(w ,map[string]string{"message":err.Error()},statusCode(err))
		return
	}
	
	mangas, err := a.Repo.GetAuthorManga(ctx,id)
	if err != nil {
		slog.Error("error at calling procedure", "message", err)
		JSONError(w, map[string]string{"message" : err.Error()}, statusCode(err))
		return
	}

	type authorManga struct{
		model.Author 
		Manga []model.Manga `json:"Manga"`
	}	

	jsonResp, err := JSONMarshaller("data succesfully retrieved",&authorManga{
		Author: *author,
		Manga: mangas,
	})

	if err != nil {
		slog.Error("error marhalling" , "error" , err)
		JSONError(w, map[string]string{"message": err.Error()}, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}	