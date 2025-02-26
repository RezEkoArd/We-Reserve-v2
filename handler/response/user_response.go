package response

import "time"

type UserResponse struct {
	ID int
	Name string
	Email string
	Role string
}


type ListUserResponse struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}