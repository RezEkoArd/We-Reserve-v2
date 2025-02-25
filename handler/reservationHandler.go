package handler

import (
	"net/http"
	"strconv"
	"time"
	"wereserve/dto"
	"wereserve/handler/response"
	"wereserve/models"
	"wereserve/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReservationHandler struct {
	ReservationService *services.ReservationService
	validate           *validator.Validate
}

func NewReservationsHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{ReservationService: reservationService,
		validate:           validator.New(),
	}
}

func (h *ReservationHandler) GetAllReservation(c *gin.Context) {
	reservations, err := h.ReservationService.GetAllReservation()
	if err != nil {
        // Tangani error khusus untuk "no reservations found"
        if err.Error() == "no reservations found" {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "No reservations found",
            })
            return
        }

        // Tangani error lainnya
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }


	// Convert to ReseponseStruct
	var allReserveResponse []response.ReservationResponse
	for _, reservation := range reservations{
		allReserveResponse = append(allReserveResponse, response.ReservationResponse{
			ID:                  reservation.ID,
			User:                response.UserPreloadResponse{
				ID:    reservation.User.ID,
				Name:  reservation.User.Name,
				Email: reservation.User.Email,
			},
			Table:               response.TableResponse{
				ID:        reservation.Table.ID,
				TableName: reservation.Table.TableName,
				Capacity:  reservation.Table.Capacity,
				Status:    reservation.Table.Status,
			},
			ReservationDateTime: reservation.ReservationDateTime,
			NumberOfPeople:      reservation.NumberOfPeople,
			CreatedAt:           reservation.CreatedAt,
			UpdatedAt:           reservation.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "successfully get all Reservation",
		"data" : allReserveResponse,
	})
}

func (h *ReservationHandler) GetReservationDetail(c *gin.Context) {
	// get Param
	ParamId := c.Param("id")
	id, err := strconv.Atoi(ParamId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "id tidak ditemukan"})
		return
	}

	// panggil service
	reservation, err := h.ReservationService.GetReservationDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Data tidak di temukan"})
		return
	}

	// convert to struct respons
	resp := response.ReservationResponse{
			ID:        reservation.ID,
			User:      response.UserPreloadResponse{
				ID:    reservation.User.ID,
				Name:  reservation.User.Name,
				Email: reservation.User.Email,
			},
			Table:         response.TableResponse{
				ID:        reservation.Table.ID,
				TableName: reservation.Table.TableName,
				Capacity:  reservation.Table.Capacity,
				Status:    reservation.Table.Status,
			},
			ReservationDateTime: reservation.ReservationDateTime,
			NumberOfPeople:      reservation.NumberOfPeople,
			CreatedAt:           reservation.CreatedAt,
			UpdatedAt:           reservation.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "data get successfully",
		"data" : resp,
	})
}

func (h *ReservationHandler) GetReservationByUserLogin(c *gin.Context) {
	// Get id user login now
	userIDInterface, exists :=  c.Get("userID")	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error" : "User Id tidak di temukan di context",
		})
		return
	}

	// Lakukan type assertion untuk mengkonversi ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Gagal konversi userID ke string",
		})
		return
	}

	// convert to int
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "userID tidak valid",
		})
		return
	}

	reservations, err := h.ReservationService.GetReservationByUserLogin(id)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H {
			"error" : "Gagal mengambil data reservasi",
		})
		return
	}

	var allReserveResponse []response.ReservationResponse
	for _, reservation := range reservations{
		allReserveResponse = append(allReserveResponse, response.ReservationResponse{
			ID:                  reservation.ID,
			User:                response.UserPreloadResponse{
				ID:    reservation.User.ID,
				Name:  reservation.User.Name,
				Email: reservation.User.Email,
			},
			Table:               response.TableResponse{
				ID:        reservation.Table.ID,
				TableName: reservation.Table.TableName,
				Capacity:  reservation.Table.Capacity,
				Status:    reservation.Table.Status,
			},
			ReservationDateTime: reservation.ReservationDateTime,
			NumberOfPeople:      reservation.NumberOfPeople,
			CreatedAt:           reservation.CreatedAt,
			UpdatedAt:           reservation.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "data get successfully",
		"data" : allReserveResponse,
	})
}



// Delete Reservation
func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error" : "Id invalid"})
		return
	}

	err = h.ReservationService.DeleteReservation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Reservation deleted Successfully" })
}

//Create 

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	// var reservation models.Reservation
	var reservationDTO dto.Reservation

	//bind dan validasi request body
	if err := c.ShouldBindJSON(&reservationDTO); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi format
	if _, err := time.Parse(time.RFC3339, reservationDTO.ReservationDateTime.Format(time.RFC3339)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid reservation date format"})
		return
	}

	// validate request
	if err := h.validate.Struct(reservationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	// Convert DTO To model
	reservationModel := models.Reservation{
		UserID: 		reservationDTO.UserID,
		TableID:            reservationDTO.TableID,
        ReservationDateTime: reservationDTO.ReservationDateTime,
        NumberOfPeople:     reservationDTO.NumberOfPeople,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
	}

	// Panggil service untuk membuat reservasi
	if err :=  h.ReservationService.CreateReservation(&reservationModel); err != nil {
		// handle error
		switch err.Error() {
		case "meja sudah di pesan" :
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		default :
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to create reservation"})
		}
		return
	}

	createReservation, err := h.ReservationService.GetReservationDetail(reservationModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to create reservation"})
		return
	}

	// format resp
	resp := response.ReservationResponse{
		ID:                  createReservation.ID,
		User:                response.UserPreloadResponse{
							ID:    createReservation.UserID,
							Name:  createReservation.User.Name,
							Email: createReservation.User.Email,
		},
		Table:               response.TableResponse{
							ID:        createReservation.TableID,
							TableName: createReservation.Table.TableName,
							Capacity:  createReservation.Table.Capacity,
							Status:    createReservation.Table.Status,
							CreatedAt: createReservation.Table.CreatedAt,
							UpdatedAt: createReservation.Table.UpdatedAt,
		},
		ReservationDateTime: createReservation.ReservationDateTime,
		NumberOfPeople:      createReservation.NumberOfPeople,
		CreatedAt:           createReservation.CreatedAt,
		UpdatedAt:           createReservation.UpdatedAt,
	}


	c.JSON(http.StatusCreated, gin.H{
		"message" : "Reservation created successfully",
		"data" : resp,
	})
}

func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid reservation ID"})
		return
	}

	// Bind and validate request body
	var updatedDTO dto.UpdateReservation
	if err := c.ShouldBindJSON(&updatedDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Validate request body using validator
	if err := h.validate.Struct(updatedDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	//Convert DTO to model.Reservation
	updateReservation := models.Reservation{
		TableID:            updatedDTO.TableID,
        UserID:             updatedDTO.UserID,
        ReservationDateTime: updatedDTO.ReservationDateTime,
        NumberOfPeople:     updatedDTO.NumberOfPeople,
	}

	//Call service to update reservation
	if err := h.ReservationService.UpdateReservation(id, updateReservation); err != nil {
		switch err.Error(){
		case "minimal 1 field harus diisi":
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		case "table_id dan reservation_datetime sudah digunakan oleh reservasi lain":
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        default:
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to update reservation"})
		}
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message" : "Reservation update successfuly",
	})
}