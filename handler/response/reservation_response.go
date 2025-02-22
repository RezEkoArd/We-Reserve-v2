package response

import "time"

type ReservationResponse struct {
    ID             int                 `json:"id"`
    User           UserPreloadResponse        `json:"user"`           // Preloaded User data
    Table          TableResponse       `json:"table"`          // Preloaded Table data
    Date           time.Time           `json:"date"`
    Time           time.Time           `json:"time"`
    NumberOfPeople int                 `json:"number_of_people"`
    CreatedAt      time.Time           `json:"created_at"`
    UpdatedAt      time.Time           `json:"updated_at"`
}

type UserPreloadResponse struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TablePreloadResponse struct {
    ID          int       `json:"id"`
    TableName   string    `json:"table_name"`
    Capacity    int       `json:"capacity"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at,omitempty"`
    UpdatedAt   time.Time `json:"updated_at,omitempty"`
}