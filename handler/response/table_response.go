package response

import "time"

type TableResponse struct {
	ID          int       `json:"id"`
    TableName string    `json:"table_name"`
    Capacity    int       `json:"capacity"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}