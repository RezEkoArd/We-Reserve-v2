package repository

import (
	"errors"
	"fmt"
	"time"
	"wereserve/models"

	"gorm.io/gorm"
)


type ReservationRepository struct {
	DB *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

func (r *ReservationRepository) GetAllReservation() ([]models.Reservation, error) {
	var reservations []models.Reservation	
	err := r.DB.Preload("User").Preload("Table").Find(&reservations)
	if err.Error != nil {
		return nil, fmt.Errorf("failed to fetch reservations : %w", err.Error)
	}

	return reservations, nil
}

func (r *ReservationRepository) GetReservationDetail(id int) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.DB.Preload("User").Preload("Table").First(&reservation, id).Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}


func (r *ReservationRepository) GetReservationByUserLogin(userID int) ([]models.Reservation, error) {
	var selfReservations []models.Reservation
	err := r.DB.Preload("User").Preload("Table").Where("user_id = ?", userID).Find(&selfReservations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reservations: %w", err)
	}

	return selfReservations, nil
}

func (r *ReservationRepository) IsReservationExists(tableID int, reservationDatetime time.Time) (bool, error) {
	var count int64
	// query untuk menghitung jumlah reservasi yang sudah ada di tableID dan reservationDatetime

	err := r.DB.Model(&models.Reservation{}).Where("table_id = ? AND reservation_datetime = ? ",tableID, reservationDatetime).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check reservation existance: %w", err)
	}

	return count > 0, nil
} 



func (r *ReservationRepository) CreateReservation(reservation *models.Reservation) error {
	reservationExists, err := r.IsReservationExists(reservation.TableID, reservation.ReservationDateTime)
	if err != nil {
		return err
	}

	if reservationExists {
		return errors.New("reservasi sudah terdaftar sebelumnya")
	}

	err = r.DB.Create(reservation).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepository) DeleteReservation(id int) error {
	sqlQuery := `DELETE FROM reservations WHERE id = $1`
	result := r.DB.Exec(sqlQuery, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("reservation with ID %d not found", id) 
	}

	return nil
}



func (r *ReservationRepository) UpdateReservation(id int, reservation *models.Reservation) error {
	var currentReservation models.Reservation
	err := r.DB.First(&currentReservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation tidak di temukan")
		}
		return err
	}

	// validasi 
	if reservation.TableID != currentReservation.TableID || !reservation.ReservationDateTime.Equal(currentReservation.ReservationDateTime) {
		exists, err := r.IsReservationExists(reservation.TableID, reservation.ReservationDateTime)
		if err != nil {
			return fmt.Errorf("gagal menvalidasi reservasi: %w", err)
		}
		if exists {
			return errors.New("table_id dan reservation sudah digunakan oleh reservasi lain")
		}
	}

	//update
	err = r.DB.Model(&currentReservation).Updates(reservation).Error
	if err != nil {
		return fmt.Errorf("gagal update reservasi : %w", err)
	}

	return nil
}