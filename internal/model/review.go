package model

type tag string

const (
	RECOMENDED     tag = "Reccomendeed"
	MIXED_FEELINGS tag = "Mixed Feelings"
	NOT_RECOMENDED tag = "Not reccomended"
)

type Review struct {
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
	User_id     string `json:"user_id,omitempty"`
	Manga_id 	string `json:"manga_id,omitempty"`
	Review_text *string `json:"review"`
	Tag         *tag    `json:"tag"`
}