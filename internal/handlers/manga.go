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

type MangaService interface {
	CreateManga(ctx context.Context, manga *model.Manga,userId int) error
	GetMangaById(ctx context.Context, id string) (*model.MangaResponse, error)
	ConnectMangaAuthor(ctx context.Context, obj *model.MangaAuthorPivot, userId int) error
	ConnectMangaGenre(ctx context.Context, obj *model.MangaGenrePivot, userId int) error
	GetAllMangaWithLimit(ctx context.Context, limit int) ([]*model.MangaList, error)
	SearchMangaByName(ctx context.Context, name string) ([]model.Manga, error)
}

type MangaHandler struct {
	Svc MangaService
}

func NewMangaHandler(mux *http.ServeMux, svc MangaService) {
	handler := &MangaHandler{
		Svc: svc,
	}

	mux.HandleFunc("GET /manga", handler.GetAllMangaWithLimit)
	mux.HandleFunc("GET /manga/{id}", handler.GetMangaById)
	mux.Handle("POST /manga", middleware.Auth(http.HandlerFunc(handler.CreateManga)))
	mux.Handle("POST /manga/{id}/author", middleware.Auth(http.HandlerFunc(handler.ConnectMangaAuthor)))
	mux.Handle("POST /manga/{id}/genre", middleware.Auth(http.HandlerFunc(handler.ConnectMangaGenre)))
}

func (m *MangaHandler) CreateManga(w http.ResponseWriter, r *http.Request) {
	var manga model.Manga

	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}	

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
		
	if err := m.Svc.CreateManga(r.Context(), &manga, userId); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Manga Created successfully"})
}

func (m *MangaHandler) GetMangaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()

	resp, err := m.Svc.GetMangaById(ctx, id)
	if err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (m *MangaHandler) ConnectMangaAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mangaId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	var payload struct {
		AuthorId int `json:"authorId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}

	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}	

	if err := m.Svc.ConnectMangaAuthor(ctx, &model.MangaAuthorPivot{Manga_id: mangaId, Author_id: payload.AuthorId}, userId); err != nil {
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Author and Manga connected successfully"})
}

func (m *MangaHandler) ConnectMangaGenre(w http.ResponseWriter, r *http.Request) {
	mangaId, err := strconv.Atoi(r.PathValue("id"))
	ctx := r.Context()
	if err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, http.StatusInternalServerError)
		return
	}
	var payload struct {
		GenreId int `json:"genreId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	if err := m.Svc.ConnectMangaGenre(ctx, &model.MangaGenrePivot{Manga_id: mangaId, Genre_id: payload.GenreId},userId); err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre and Manga connected successfully"})
}

func (m *MangaHandler) GetAllMangaWithLimit(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("title")
	if query != ""{
		ctx := r.Context()
		mangas, err := m.Svc.SearchMangaByName(ctx, query)
		if err != nil {
			slog.Error("error at executing query","message",err)
			JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
			return
		}

		jsonResp, err := json.Marshal(mangas)
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
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		JSONError(w,  map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	mangaList, err := m.Svc.GetAllMangaWithLimit(ctx, limit)
	if err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return
	}

	jsonResp, err := json.Marshal(mangaList)
	if err != nil {
		slog.Error("error marshal json")
		JSONError(w,  map[string]string{"message": err.Error()}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}