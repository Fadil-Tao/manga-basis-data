package model

type status string

const (
	FINISHED    status = "finished"
	IN_PROGRESS status = "in_progress"
)

type Manga struct {
	Id           string  `json:"id"`
	Title        string  `json:"title"`
	Synopsys     string  `json:"synopsys,omitempty"`
	Manga_status status  `json:"status"`
	Published_at string  `json:"published_at"`
	Finished_at  *string `json:"finished_at,omitempty"`
}

type MangaAuthorPivot struct {
	Manga_id  int
	Author_id int
}

type MangaGenrePivot struct {
	Manga_id int
	Genre_id int
}

type MangaResponse struct {
	Manga
	Genres []Genre  `json:"genre"`
	Author []Author `json:"author"`
}



type MangaList struct{
	Manga	
	Rating float64 	`json:"rating"`
	TotalReview float64 `json:"totalReview"`
	TotalLikes float64 `json:"likes"`
	TotalUserRated float64 `json:"totalUserRated"`
}