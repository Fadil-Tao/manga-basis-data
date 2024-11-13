package model

type Rating struct {	
	Id string `json:"id"`
	Manga_id string `json:"mangaiId"`
	User_id string `json:"userId"`
	Created_at string `json:"createdAt"` 
}