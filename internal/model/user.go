package model


type User struct{
	Id int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string
	Created_at string `json:"creted_at"`
	Updated_at string 
	Is_admin int8 `json:"Is_admin"`	
	Is_deleted int8 `json:"Is_deleted"`	
}	
type UserResponse struct{
	Id int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Is_admin int `json:"is_admin"`
	Created_at string `json:"created_at"`
}
type NewUserRequest struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct{
	Email string `json:"email"`
	Password  string `json:"password"`
}