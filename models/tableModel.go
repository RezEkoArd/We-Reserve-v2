package models

import "time"

type Table struct {
	ID          int       `json:"id"`
    TableNumber string    `json:"table_number"`
    Capacity    int       `json:"capacity"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

