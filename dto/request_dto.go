package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

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


// Table Validator

type CreateTableRequest struct {
	TableName string `json:"table_name" validate:"required,min=3,max=50"`
	Capacity  int    `json:"capacity" validate:"required,min=1,max=20"`
	Status    string `json:"status" validate:"required,oneof=available reserved occupied"`
}

type UpdateTableRequest struct {
	TableName string `json:"table_name" validate:"omitempty,min=3,max=50"`
	Capacity  int    `json:"capacity" validate:"omitempty,min=1,max=20"`
	Status    string `json:"status" validate:"omitempty,oneof=available reserved occupied"`
}


// Reservation Validator
type UpdateReservationRequest struct {
	Date *time.Time `json:"date" validate:"omitempty"`
	Time *time.Time	`json:"time" validate:"omitempty"`
	NumberOfPeople *int	`json:"number_of_people" validate:"omitempty,min=1"`
}


type CreateReservationRequest struct {
	UserID int	`json:"user_id" validate:"required"`
	TableID	int	`json:"table_id" validate:"required"`
	Date	string `json:"date" validate:"required"`
	Time	string `json:"time" validate:"required"`
	NumberOfPeople int	`json:"number_of_people" validate:"required,min=1"`
}

// Validator Instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}
