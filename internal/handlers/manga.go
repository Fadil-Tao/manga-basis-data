package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type MangaService interface{
	CreateManga(ctx context.Context, manga *model.Manga)error
	GetMangaById(ctx context.Context, id string) (*model.MangaResponse, error)
	ConnectMangaAuthor(ctx context.Context, obj *model.MangaAuthorPivot) error
	ConnectMangaGenre(ctx context.Context, obj *model.MangaGenrePivot) error
	GetAllMangaWithLimit(ctx context.Context,limit int)([]*model.MangaList, error)
}

type MangaHandler struct {
	Svc MangaService
}

func NewMangaHandler(mux *http.ServeMux,svc MangaService){
	handler := &MangaHandler{
		Svc: svc,
	}	
	
	mux.HandleFunc("GET /manga", handler.GetAllMangaWithLimit)
	mux.HandleFunc("GET /manga/{id}", handler.GetMangaById)
	mux.HandleFunc("POST /manga", handler.CreateManga)
	mux.HandleFunc("POST /manga/{id}/author", handler.ConnectMangaAuthor)
	mux.HandleFunc("POST /manga/{id}/genre", handler.ConnectMangaGenre)
}

func (m *MangaHandler)CreateManga(w http.ResponseWriter, r *http.Request){
	var manga model.Manga
	
	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Svc.CreateManga(r.Context(), &manga); err != nil {
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Manga Created successfully"})
}

func (m *MangaHandler) GetMangaById(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")

	ctx := r.Context()

	resp , err := m.Svc.GetMangaById(ctx, id)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (m *MangaHandler) ConnectMangaAuthor(w http.ResponseWriter , r *http.Request){
	ctx := r.Context()
	mangaId,err := strconv.Atoi(r.PathValue("id"))
	if err !=nil {
		http.Error(w, "failed convert id to int",http.StatusInternalServerError)
		return
	}
	var payload struct{
		AuthorId int `json:"authorId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
	if err := m.Svc.ConnectMangaAuthor(ctx, &model.MangaAuthorPivot{Manga_id: mangaId, Author_id: payload.AuthorId}); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author and Manga connected successfully"})	
} 

func (m *MangaHandler) ConnectMangaGenre(w http.ResponseWriter , r *http.Request){
	mangaId,err := strconv.Atoi(r.PathValue("id"))
	ctx := r.Context()
	if err != nil{
		http.Error(w, "failed convert id to int",http.StatusInternalServerError)
		return	
	}
	var payload struct{
		GenreId int `json:"genreId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
	if err := m.Svc.ConnectMangaGenre(ctx, &model.MangaGenrePivot{Manga_id: mangaId, Genre_id: payload.GenreId}); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre and Manga connected successfully"})	
} 


func (m *MangaHandler) GetAllMangaWithLimit(w http.ResponseWriter, r *http.Request){
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10" 
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	mangaList,err := m.Svc.GetAllMangaWithLimit(ctx,limit)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(mangaList)
	if err != nil {
		slog.Error("error marshal json")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}




