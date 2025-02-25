package models

import "time"

type Reservation struct {
	ID             int       `json:"id"`
    UserID         int       `json:"user_id" gorm:"column:user_id"`
    TableID        int       `json:"table_id" gorm:"column:table_id"`
    ReservationDateTime          time.Time    `json:"reservation_datetime" gorm:"column:reservation_datetime"`
    NumberOfPeople int       `json:"number_of_people" gorm:"column:number_of_people"`
    CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`

    // Relasi model ke User dan ke Table
    User User `gorm:"foreignKey:UserID"`
    Table Table `gorm:"foreignKey:TableID"`
}