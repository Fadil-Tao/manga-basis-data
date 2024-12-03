package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetAllManga(ctx context.Context, name string) ([]*model.MangaList, error) 
	SearchMangaByName(ctx context.Context, title string) ([]model.Manga, error)
	DeleteMangaById(ctx context.Context, id int, userId int) error 
	UpdateManga(ctx context.Context,id int,manga *model.Manga, userId int)error
	GetMangaRankingList(ctx context.Context, period string)([]*model.MangaList, error)
	ToggleLikeManga(ctx context.Context, userId int , mangaId int)error
	DeleteMangaAuthorConnection(ctx context.Context, userId int, mangaId int, authorId int)error
	DeleteMangaGenreConnection(ctx context.Context, userId int, mangaId int, genreId int)error
}

type MangaHandler struct {
	Svc MangaService
}

func NewMangaHandler(mux *http.ServeMux, svc MangaService) {
	handler := &MangaHandler{
		Svc: svc,
	}

	mux.HandleFunc("GET /manga", handler.GetAllManga)
	mux.HandleFunc("GET /manga/{id}", handler.GetMangaById)
	mux.HandleFunc("GET /manga/rank", handler.GetMangaRankList)
	mux.Handle("POST /manga", middleware.Auth(http.HandlerFunc(handler.CreateManga)))
	mux.Handle("POST /manga/{id}/author", middleware.Auth(http.HandlerFunc(handler.ConnectMangaAuthor)))
	mux.Handle("POST /manga/{id}/genre", middleware.Auth(http.HandlerFunc(handler.ConnectMangaGenre)))
	mux.Handle("DELETE /manga/{id}" , middleware.Auth(http.HandlerFunc(handler.DeleteMangaById)))
	mux.Handle("PUT /manga/{id}", middleware.Auth(http.HandlerFunc(handler.UpdateManga)))
	mux.Handle("POST /manga/{id}/like", middleware.Auth(http.HandlerFunc(handler.ToggleLikeManga)))
	mux.Handle("DELETE /manga/{mangaId}/author/{authorId}", middleware.Auth(http.HandlerFunc(handler.DeleteMangaAuthorConnection)))
	mux.Handle("DELETE /manga/{mangaId}/genre/{genreId}", middleware.Auth(http.HandlerFunc(handler.DeleteMangaGenreConnection)))
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

	jsonResp, err := JSONMarshaller("manga success retrieved",resp)
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

func (m *MangaHandler) GetAllManga(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var queryMap = map[string]string{
		"name":"",
	}
	if name := query.Get("name"); name != ""{
		queryMap["name"] = name
	}
	validPeriods := map[string]bool{
		"today" : true,
		"month": true,
		"all": true,
	}
	rank := query.Get("rank");
	if validPeriods[rank]{
		mangas, err := m.Svc.GetMangaRankingList(r.Context(), rank)
		if err != nil {
			if statusCode(err) != http.StatusInternalServerError{
				JSONError(w, map[string]string{"message": err.Error()} , statusCode(err))
				return
			}
			JSONError(w, map[string]string{"message": "internal server error"} , http.StatusInternalServerError)
			return 
		}
		jsonResp, err := JSONMarshaller("succefully retrieved manga", mangas)
		if err != nil {
			JSONError(w, map[string]string{"message": "internal server error"}, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return
	}
	mangas, err := m.Svc.GetAllManga(r.Context(), queryMap["name"])
	if err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{"message": err.Error()} , statusCode(err))
			return
		}
		JSONError(w, map[string]string{"message": "internal server error"} , http.StatusInternalServerError)
		return 
	}
	jsonResp, err := JSONMarshaller("succefully retrieved manga", mangas)
	if err != nil {
		JSONError(w, map[string]string{"message": "internal server error"}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}


func(m *MangaHandler) DeleteMangaById(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
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

	if err := m.Svc.DeleteMangaById(ctx, id , userId); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Manga successfully deleted"})
}

func (m *MangaHandler)UpdateManga(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	var manga model.Manga
	if err := json.NewDecoder(r.Body).Decode(&manga);err != nil{
		JSONError(w,map[string]string{"message":err.Error()}, http.StatusBadRequest)
		return
	}
	
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return 
	}

	if err := m.Svc.UpdateManga(ctx, id, &manga, userId); err != nil {
		slog.Error("error using repo", "message", err)
		JSONError(w, map[string]string{"message":err.Error()}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Manga updated successfully"})
}

func (m *MangaHandler) GetMangaRankList(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("period")
	if query != ""{
		ctx := r.Context()
		mangas, err := m.Svc.GetMangaRankingList(ctx,query)
		if err != nil{
			slog.Error("error at processing param", "message", err)
			JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
			return
		}
		jsonResp, err := JSONMarshaller("manga rank list successfully retrieved",mangas)
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
	ctx := r.Context()
	mangas, err := m.Svc.GetMangaRankingList(ctx,"all")
	if err != nil{
		slog.Error("error at processing param", "message", err)
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
}


func (m *MangaHandler) ToggleLikeManga(w http.ResponseWriter, r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return 
	}
	mangaId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":"Internal server error"}, http.StatusInternalServerError)
		return
	}

	if err := m.Svc.ToggleLikeManga(r.Context() , userId,mangaId); err != nil {
		slog.Error("error exec procedure" , "error", err)
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Toggle triggered successfully"})
}
func(m *MangaHandler) DeleteMangaAuthorConnection(w http.ResponseWriter,  r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return 
	}
	mangaId, err := strconv.Atoi(r.PathValue("mangaId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":"Internal server error"}, http.StatusInternalServerError)
		return
	}
	authorId, err := strconv.Atoi(r.PathValue("authorId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":"Internal server error"}, http.StatusInternalServerError)
		return
	}
	if err := m.Svc.DeleteMangaAuthorConnection(r.Context(), userId,mangaId,authorId); err != nil {
		slog.Info("author =id",":",authorId)
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}
func(m *MangaHandler) DeleteMangaGenreConnection(w http.ResponseWriter,  r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return 
	}
	mangaId, err := strconv.Atoi(r.PathValue("mangaId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":"Internal server error"}, http.StatusInternalServerError)
		return
	}
	genreId, err := strconv.Atoi(r.PathValue("genreId"))
	if err != nil {
		slog.Error("error converting id to int")
		JSONError(w, map[string]string{"message":"Internal server error"}, http.StatusInternalServerError)
		return
	}
	fmt.Print(genreId)
	if err := m.Svc.DeleteMangaGenreConnection(r.Context(), userId,mangaId,genreId); err != nil {
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}