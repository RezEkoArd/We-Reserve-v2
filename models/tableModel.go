package models

import "time"

type Table struct {
	ID          int       `json:"id"`
    TableName string    `json:"table_name" gorm:"column:table_name"`
    Capacity    int       `json:"capacity"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

