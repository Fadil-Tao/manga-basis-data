package model

type Manga struct{
	Id string  `json:"id"`	
	Title string `json:"title"`
	Synopsys string `json:"synopsys"`
	Published_at string `json:"published_at"`
}

