package dto

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Name string `json:"name" validate:"required,min=3,max=20"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role string `json:"role" validate:"omitempty,oneof=admin customer"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=3,max=20"`
	Email string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}


// Validator Instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}