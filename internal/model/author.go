package model


type Author struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Birthday string `json:"birthday,omitempty"`
	Biography string `json:"biography,omitempty"`
}
