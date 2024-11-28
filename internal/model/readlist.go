package model 

type Readlist struct {
	Id string `json:"id"`
	UserId string `json:"userId,omitempty"`
	UserName string `json:"owner"`
	Name string `json:"name"`
	Description string `json:"description"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}


type NewReadlistRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type ReadlistItem struct {
	Id string `json:"id"`
	MangaId string `json:"mangaId"`
	Title string `json:"title"`
	Status string `json:"readStatus"`
	AddedAt string `json:"addedAt"`
} 

type NewReadlistItem struct {
	MangaId string `json:"mangaId"`
	Status string `json:"readStatus"`
	ReadlistId string `json:"readListId"`
}