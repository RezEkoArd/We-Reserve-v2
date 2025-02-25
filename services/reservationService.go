package services

import (
	"errors"
	"fmt"
	"wereserve/models"
	"wereserve/repository"

	"github.com/go-playground/validator/v10"
)


type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	Validator *validator.Validate
	tableRepo *repository.TableRepository
}

func NewReservationService(reservationRepo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{reservationRepo:  reservationRepo}
}


//Get All Reservation 
func (s *ReservationService) GetAllReservation() ([]models.Reservation, error) {
	reservations, err := s.reservationRepo.GetAllReservation()
	if err != nil {
		return nil, err
	}

	return reservations, nil
}



func (s *ReservationService) GetReservationDetail(id int) (*models.Reservation, error) {
	
	// validasi id tersedia di db atau tidak
	reservation, err := s.reservationRepo.GetReservationDetail(id)
	if err != nil {
		return nil,err
	}

	return reservation, nil
}

func (s *ReservationService) GetReservationByUserLogin(userID int) ([]models.Reservation, error) {
	reservations, err := s.reservationRepo.GetReservationByUserLogin(userID)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	return reservations, nil
}

func (s *ReservationService) CreateReservation(reservation *models.Reservation) error {
	
	// Cek Reservation 
	// _, err := s.reservationRepo.IsReservationExists(reservation.TableID, reservation.ReservationDateTime)
	// if err != nil {
	// 	return fmt.Errorf("reservasi sudah dibuat: %w" , err)
	// }

	table, err := s.tableRepo.GetTableByID(reservation.TableID)
	if err != nil {
		return fmt.Errorf("gagal mengambil data meja: %w", err)
	}

	if table.Status == "reserved" {
		return errors.New("meja sudah di pesan")
	}

	// buat reservasi
	err = s.reservationRepo.CreateReservation(reservation) 
	if err != nil {
		return fmt.Errorf("gagal membuat reservasi: %w", err)
	}

	table.Status = "reserved"
	err = s.tableRepo.UpdateTable(reservation.TableID, table)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status meja: %w", err)
	}

	return nil
}

// Delete 
func (s *ReservationService) DeleteReservation(id int) error {
	_, err := s.reservationRepo.GetReservationDetail(id)
	if err != nil {
		return errors.New("reservation id tidak ditemukan")
	}

	err = s.reservationRepo.DeleteReservation(id)
	if err != nil {
		return err
	}

	return nil
}

// update 
func (s *ReservationService) UpdateReservation(id int, reservation models.Reservation) error {
	if reservation.TableID == 0 && reservation.UserID == 0 && reservation.ReservationDateTime.IsZero() && reservation.NumberOfPeople == 0 {
		return errors.New("minimal 1 field harus di isi")
	 }	
	
	// validasi reservasi yang akan diupdate ada
	 currentReservation, err := s.reservationRepo.GetReservationDetail(id)
	 if err != nil {
		return fmt.Errorf("gagal mengambil reservasi; %w", err)
	 }

	 if reservation.TableID != currentReservation.TableID || !reservation.ReservationDateTime.Equal(currentReservation.ReservationDateTime) {
		exists, err := s.reservationRepo.IsReservationExists(reservation.TableID, reservation.ReservationDateTime)
		if err != nil {
			return fmt.Errorf("gagal validasi reervasi: %w", err)
		}
		if exists {
			return errors.New("table_id dan reservation_datetime sudah digunakan oleh reservasi lain")
		}
	 }


	 // update reservasi
	 err = s.reservationRepo.UpdateReservation(id, &reservation)
	 if err != nil {
		return fmt.Errorf("gagal mengupdate reservasi: %w", err)
	 }

	 return nil
}