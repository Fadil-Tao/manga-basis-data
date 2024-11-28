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

type ReadlistRepo interface {
	CreateReadlist(ctx context.Context, userId int,readlist *model.NewReadlistRequest)error
	DeleteReadlist(ctx context.Context, id int, userId int)error
	UpdateReadlist(ctx context.Context, id int, userId int,readlist *model.NewReadlistRequest) error
	SearchReadlist(ctx context.Context, name string) ([]*model.Readlist, error)
	GetUserReadlist(ctx context.Context, userId int) ([]*model.Readlist, error)
	AddToReadlist(ctx context.Context, userId int,readlistItem *model.NewReadlistItem)error
	DeleteReadlistItem(ctx context.Context, userId int, readlistItemId int) error
	UpdateReadlistItemStatus(ctx context.Context, readlistId int,userId int,status string) error 
	GetReadlistItem(ctx context.Context, id int)([]*model.ReadlistItem, error)
}

type ReadlistHandler struct {
	Repo ReadlistRepo
}


func NewReadlistHandler(mux *http.ServeMux,repo ReadlistRepo){
	handler := &ReadlistHandler{
		Repo: repo,
	}

	mux.HandleFunc("GET /readlist", handler.SearchReadlist)                      
	mux.HandleFunc("GET /user/{id}/readlist", handler.GetUserReadlist)     
	mux.HandleFunc("GET /readlist/{id}/item", handler.GetReadlistItemFromReadlist)     
	mux.Handle("POST /readlist", middleware.Auth(http.HandlerFunc(handler.CreateReadlist))) 
	mux.Handle("DELETE /readlist/{id}", middleware.Auth(http.HandlerFunc(handler.DeleteReadlist))) 
	mux.Handle("PUT /readlist/{id}", middleware.Auth(http.HandlerFunc(handler.UpdateReadlist)))  

	mux.Handle("POST /readlist/item", middleware.Auth(http.HandlerFunc(handler.AddToReadlist))) 
	mux.Handle("DELETE /readlist/item/{id}", middleware.Auth(http.HandlerFunc(handler.DeleteReadlistItem))) 
	mux.Handle("PATCH /readlist/item/{id}", middleware.Auth(http.HandlerFunc(handler.UpdateReadlistItemStatus)))
}


func (rl *ReadlistHandler) CreateReadlist(w http.ResponseWriter, r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	var Readlist model.NewReadlistRequest 
	if err := json.NewDecoder(r.Body).Decode(&Readlist); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := rl.Repo.CreateReadlist(r.Context(), userId,&Readlist) ; err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Readlist Created successfully"})
}


func (rl *ReadlistHandler) DeleteReadlist(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	
	if err := rl.Repo.DeleteReadlist(r.Context(),id , userId); err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Readlist deleted successfully"})	
}


func (rl *ReadlistHandler) UpdateReadlist(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	var Readlist model.NewReadlistRequest
	if err := json.NewDecoder(r.Body).Decode(&Readlist); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := rl.Repo.UpdateReadlist(r.Context(), id, userId,&Readlist); err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Readlist updated successfully"})	
}

func (rl *ReadlistHandler) SearchReadlist(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("name")
	name := ""
	
	if query != ""{
		name = query
	}
	readlists , err := rl.Repo.SearchReadlist(r.Context(), name)
	if err != nil {
		slog.Error("error at calling procedure", "message", err)
		JSONError(w, map[string]string{"message": err.Error()}, statusCode(err))
		return	
	}
	jsonResp, err := JSONMarshaller("readlist succesfully retrieved",readlists)
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

func (rl *ReadlistHandler) GetUserReadlist(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	
	readlists,err := rl.Repo.GetUserReadlist(r.Context(),id);
	if  err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	jsonResp, err := JSONMarshaller("readlist succesfully retrieved", readlists)
	if err != nil {
		JSONError(w, map[string]string{
			"message":"internal server error",
		}, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}	


func (rl *ReadlistHandler) AddToReadlist(w http.ResponseWriter, r *http.Request){
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	var manga model.NewReadlistItem 
	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := rl.Repo.AddToReadlist(r.Context(),userId,&manga); err != nil{
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "manga added to readlist successfully"})
}

func (rl *ReadlistHandler) GetReadlistItemFromReadlist(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	mangas,err := rl.Repo.GetReadlistItem(r.Context(), id);
	if err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return	
	}
	jsonResp, err := JSONMarshaller("readlist item succesfully retrieved", mangas)
	if err != nil {
		JSONError(w, map[string]string{
			"message":"internal server error",
		}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (rl *ReadlistHandler) DeleteReadlistItem(w  http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	if err := rl.Repo.DeleteReadlistItem(r.Context(),userId, id); err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Readlist item deleted successfully"})		
}

func (rl *ReadlistHandler) UpdateReadlistItemStatus(w http.ResponseWriter, r *http.Request){
	id , err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}
	userId, err := middleware.GetUserId(w, r)
	if err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	statusReadlist := struct {
		Status string `json:"status"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&statusReadlist); err != nil {
		JSONError(w, map[string]string{
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}


	if err := rl.Repo.UpdateReadlistItemStatus(r.Context(),id,userId,statusReadlist.Status); err != nil {
		if statusCode(err) != http.StatusInternalServerError{
			JSONError(w, map[string]string{
				"message": err.Error(),
			}, statusCode(err))
			return
		}
		JSONError(w, map[string]string{
			"message": "internal server error",
		}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Readlist updated successfully"})	
}

