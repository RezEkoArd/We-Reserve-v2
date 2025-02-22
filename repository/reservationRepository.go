package repository

import (
	"errors"
	"fmt"
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
	result := r.DB.Preload("User").Preload("Table").Find(&reservations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fecth reservations: %w", result.Error)
	}

	return reservations, nil
}

//Get Reservation by using account now
func (r *ReservationRepository) GetReservationsByUsers(userEmail string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	result := r.DB.Preload("User").Preload("Table").Where("email = ?", userEmail).Find(&reservations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch reservation for user %s: %w", userEmail, result.Error)
	}
	return reservations, nil
}


//Delete Reservation
func (r *ReservationRepository) DeleteReservation(id int) error {
	sql_query := `DELETE FROM reservations WHERE id = $1`
	result := r.DB.Exec(sql_query, id)

	if result.Error != nil {
		return result.Error
	}

	// validate if err 
	if result.RowsAffected == 0 {
		return fmt.Errorf("reservation is not found")
	}

	return nil
}

func (r *ReservationRepository) CreateReservation(reservation *models.Reservation) error {
	//Validasi apakah meja sudah dipesan pada tgl tertentu
	var existingReservation models.Reservation
	result := r.DB.Where("table_id = ? AND date = ? AND time = ?",
	reservation.TableID, reservation.Date, reservation.Time).First(&existingReservation)

	if result.Error == nil {
		//jika reservation sudah ada, kembalikan error
		return fmt.Errorf("table %d is already reserve for %s at %s",
		reservation.TableID,reservation.Date, reservation.Time)
	}


	// buat Reservation
	result = r.DB.Create(reservation)
	if result.Error != nil {
		// Jika gagal membuat reservasi, kembalikan error
		return fmt.Errorf("failed to created reservation: %w", result.Error)
	}


	return nil
}


func (r *ReservationRepository) UpdateReservation(id int, updatedReservation *models.Reservation) error {
	var reservation *models.Reservation
	err := r.DB.First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation tidak di temukan")
		}
		return err
	}
	
	var existingReservation models.Reservation
	result := r.DB.Where("table_id = ? AND date = ? AND time = ?",
	reservation.TableID, reservation.Date, reservation.Time).First(&existingReservation)

	if result.Error == nil {
		//jika reservation sudah ada, kembalikan error
		return fmt.Errorf("table %d is already reserve for %s at %s",
		reservation.TableID,reservation.Date, reservation.Time)
	}

	// Perbarui data reservasi
	updateResult := r.DB.Model(&reservation).Updates(updatedReservation)
	if result.Error != nil {
		return fmt.Errorf("failed to update reservation: %w", updateResult.Error)
	} 

	return nil
}	


func (r *ReservationRepository) GetReservationDetail(id int) (*models.Reservation, error) {
	var reservation *models.Reservation
	err := r.DB.Raw("SELECT * FROM reservations WHERE id = $1", id).First(&reservation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("reservation not found")
		}
		return reservation, nil
	}

	return reservation, nil
}