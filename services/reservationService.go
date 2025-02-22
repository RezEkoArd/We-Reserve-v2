package services

import (
	"errors"
	"fmt"
	"time"
	"wereserve/dto"
	"wereserve/models"
	"wereserve/repository"

	"github.com/go-playground/validator/v10"
)


type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	Validator *validator.Validate
}

func NewReservationService(reservationRepo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{reservationRepo:  reservationRepo}
}

func (s *ReservationService) CreateReservation(userID, tableID int, dateStr, timeStr string, numberOfPeople int) error {
	// Validasi jumlah orang
	if numberOfPeople <= 0 {
		return fmt.Errorf("number of people must be greater tha zero")
	}

	// parse date string to time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	fmt.Println("Parsed Date:", date) 

	ParsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format: %w" , err)
	}
	fmt.Println("Parsed Time:", ParsedTime) 

	// combine date and time 
	// datetime := time.Date(date.Year(),date.Month(),date.Day(),ParsedTime.Hour(),ParsedTime.Minute(),0,0,time.Local) 

	// Buat object reservation
	reservation := &models.Reservation{
		UserID:         userID,
		TableID:        tableID,
		Date:           date,
		Time:           ParsedTime,
		NumberOfPeople: numberOfPeople,
	}

	//kirim ke repository
	err = s.reservationRepo.CreateReservation(reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation : %w", err)
	}

	return nil
}

func (s *ReservationService) GetAllReservation() ([]models.Reservation, error) {
	reservations, err := s.reservationRepo.GetAllReservation()
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (s *ReservationService) GetReservationByUsers(userEmail string) ([]models.Reservation, error) {
	reservation, err := s.reservationRepo.GetReservationsByUsers(userEmail)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) GetReservationDetail(id int) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.GetReservationDetail(id) 
	if err != nil {
		return nil, err
	}

	return reservation, nil
}


func (s *ReservationService) DeleteReservation(id int) error {
	// Cek Id Reservation
	_, err := s.reservationRepo.GetReservationDetail(id) 
	if err != nil {
		return errors.New("id reservation tidak di temukan")
	}	

	// Delete
	err = s.reservationRepo.DeleteReservation(id)
	if err != nil {
		return err
	}
	return nil
}


func (s *ReservationService) UpdateReservation(id int, req dto.UpdateReservationRequest) error {
	// validasi ID
	reservation, err := s.reservationRepo.GetReservationDetail(id)
	if err != nil {
		return errors.New("reservation not found")
	}

	// step 2: validasi Request Body
	if err := s.Validator.Struct(req); err != nil {
		return err
	}

	// Pastikan minimal 1 vield
	if req.Date == nil && req.Time == nil && req.NumberOfPeople == nil {
		return errors.New("at least one field must be provided for update")
	}

	// Cek apakah ada konf reservasi
	var updatedDate time.Time
	var	updatedTime time.Time
	var updatedNumberOfPeople int

	if req.Date != nil {
        updatedDate = *req.Date
    } else {
        updatedDate = reservation.Date
    }

    if req.Time != nil {
        updatedTime = *req.Time
    } else {
        updatedTime = reservation.Time
    }

    if req.NumberOfPeople != nil {
        updatedNumberOfPeople = *req.NumberOfPeople
    } else {
        updatedNumberOfPeople = reservation.NumberOfPeople
    }

	// Cek apakah meja sudah dipesan pada tanggal dan waktu yang baru
	var existingReservation models.Reservation
	result := s.reservationRepo.DB.Where("table_id = ? AND date = ? AND time = ? AND id <> ?",
		reservation.TableID, updatedDate, updatedTime, id).First(&existingReservation)
		if result.Error == nil {
			return errors.New("table is already reserved for the given date and time")
		}

	updatedReservation := models.Reservation{
		Date: updatedDate,
		Time: updatedTime,
		NumberOfPeople: updatedNumberOfPeople,
	}

	if err := s.reservationRepo.UpdateReservation(id, &updatedReservation); err != nil {
		return err
	}

	return nil
}