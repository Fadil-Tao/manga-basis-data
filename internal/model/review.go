package model

type tag string

const (
	RECOMENDED     tag = "Reccomendeed"
	MIXED_FEELINGS tag = "Mixed Feelings"
	NOT_RECOMENDED tag = "Not reccomended"
)

type Review struct {
	Id          string `json:"id"`
	Manga_id    string `json:"manga_id,omitempty"`
	Username string `json:"username,omitempty"`
	User_id     string `json:"user_id"`
	Review_text string `json:"review"`
	Tag         tag    `json:"tag"`
	Created_at  string `json:"created_at,omitempty"`
	Total_Like int `json:"like"`
}

type NewReviewRequest struct {
	Manga_id    string `json:"manga_id"`
	User_id     string `json:"user_id"`
	Review_text string `json:"review"`
	Tag         tag    `json:"tag"`
}

type UpdateReview struct {
	Id          string `json:"id"`
	User_id     string `json:"user_id"`
	Review_text string `json:"review"`
	Tag         tag    `json:"tag"`
}