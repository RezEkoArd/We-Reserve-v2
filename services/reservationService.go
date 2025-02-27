package services

import (
	"errors"
	"fmt"
	"log"
	"wereserve/models"
	"wereserve/repository"
	"wereserve/utils"

	"github.com/go-playground/validator/v10"
)


type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	Validator *validator.Validate
	tableRepo *repository.TableRepository
}

func NewReservationService(reservationRepo *repository.ReservationRepository, tableRepo *repository.TableRepository) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		tableRepo: tableRepo,
        Validator:       validator.New(),
	}
}


//Get All Reservation 
func (s *ReservationService) GetAllReservation() ([]models.Reservation, error) {
	reservations, err := s.reservationRepo.GetAllReservation()
	if err != nil {
		return nil, err
	}

	// Cek apakah reservation kosong
	if len(reservations) == 0 {
		return nil, errors.New("no reservations found")
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

func (s *ReservationService) CreateReservation(reservation *models.Reservation, emailUser string) error {
	
// Validasi input menggunakan Validator
	if err := s.Validator.Struct(reservation); err != nil {
		log.Printf("Validation error : %v", err)
		return fmt.Errorf("invalid reservation data: %w", err)
	}

	// Cek apakah meja tersedia
	table, err := s.tableRepo.GetTableByID(reservation.TableID)
	if err != nil {
		log.Printf("Failed to fetch table with ID %d: %v", reservation.TableID, err)
		return fmt.Errorf("failed to fetch table with ID %d: %w", reservation.TableID, err)
	}
	
	if table.Status == "reserved" {
		log.Printf("Table with ID %d is already reserved", reservation.TableID)
		return fmt.Errorf("table with ID %d is already reserved", reservation.TableID)
	}

	// Buat reservasi
	if err := s.reservationRepo.CreateReservation(reservation); err != nil {
		log.Printf("Failed to create reservation for table ID %d: %v", reservation.TableID, err)
		return fmt.Errorf("failed to create reservation for table ID %d: %w", reservation.TableID, err)
	}

	// Update status meja menjadi "reserved"
	table.Status = "reserved"
	if err := s.tableRepo.UpdateTable(reservation.TableID, table); err != nil {
		log.Printf("Failed to update status for table ID %d: %v", reservation.TableID, err)
		return fmt.Errorf("failed to update status for table ID %d: %w", reservation.TableID, err)
	}

	err = utils.SendEmail( emailUser, table.TableName, reservation.ReservationDateTime.Format("2006-01-02 15:04"))
	if err != nil {
		log.Printf("Failed to send Email Confirmation")
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
func (s *ReservationService) UpdateReservation(id int, updatedReservation models.Reservation) error {
	// Ambil reservasi saat ini
	currentReservation, err := s.reservationRepo.GetReservationDetail(id)
	if err != nil {
		return fmt.Errorf("failed to fetch reservation with ID %d: %w", id, err)
	}

	// Periksa apakah ada perubahan
	if updatedReservation.TableID == 0 && updatedReservation.UserID == 0 &&
		updatedReservation.ReservationDateTime.IsZero() && updatedReservation.NumberOfPeople == 0 {
		return errors.New("at least one field must be updated")
	}

	// Periksa konflik reservasi
	if updatedReservation.TableID != currentReservation.TableID ||
		!updatedReservation.ReservationDateTime.Equal(currentReservation.ReservationDateTime) {
		exists, err := s.reservationRepo.IsReservationExists(updatedReservation.TableID, updatedReservation.ReservationDateTime)
		if err != nil {
			return fmt.Errorf("failed to validate reservation: %w", err)
		}
		if exists {
			return errors.New("reservation conflict: table and datetime already reserved")
		}
	}

	// Update reservasi
	if err := s.reservationRepo.UpdateReservation(id, &updatedReservation); err != nil {
		return fmt.Errorf("failed to update reservation with ID %d: %w", id, err)
	}

 return nil
}