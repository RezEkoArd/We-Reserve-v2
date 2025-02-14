package models

import "time"

type Reservation struct {
	ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    TableID        int       `json:"table_id"`
    Date           string    `json:"date"`
    Time           string    `json:"time"`
    NumberOfPeople int       `json:"number_of_people"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`

}