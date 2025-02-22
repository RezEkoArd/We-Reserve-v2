package models

import "time"

type Reservation struct {
	ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    TableID        int       `json:"table_id"`
    Date           time.Time    `json:"date"`
    Time           time.Time    `json:"time"`
    NumberOfPeople int       `json:"number_of_people"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`

    // Relasi model ke User dan ke Table
    User User `gorm:"foreignKey:UserID"`
    Table Table `gorm:"foreignKey:TableID"`
}