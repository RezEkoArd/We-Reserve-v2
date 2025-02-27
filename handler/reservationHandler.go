package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
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


// GetAllReservation godoc
// @Summary      Get all reservations
// @Description  Retrieve a list of all reservations
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Success      200  {array}   response.ReservationResponse "List of reservations retrieved successfully"
// @Failure      404  {object}  response.ErrorResponse       "No reservations found"
// @Failure      500  {object}  response.ErrorResponse       "Internal server error"
// @Router       /api/reservation [get]
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


// GetReservationDetail godoc
// @Summary      Get a reservation by ID
// @Description  Retrieve a reservation's details based on its ID
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id   path      int                      true  "Reservation ID"
// @Success      200  {object}  response.ReservationResponse "Reservation retrieved successfully"
// @Failure      400  {object}  response.ErrorResponse       "Invalid reservation ID"
// @Failure      404  {object}  response.ErrorResponse       "Reservation not found"
// @Failure      500  {object}  response.ErrorResponse       "Internal server error"
// @Router       /api/reservation/{id} [get]
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


// GetReservationByUserLogin godoc
// @Summary      Get reservations by logged-in user
// @Description  Retrieve a list of reservations made by the currently logged-in user
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Success      200  {array}   response.ReservationResponse "List of reservations retrieved successfully"
// @Failure      400  {object}  response.ErrorResponse       "Invalid user ID or failed to retrieve data"
// @Failure      401  {object}  response.ErrorResponse       "Unauthorized, user ID not found in context"
// @Failure      500  {object}  response.ErrorResponse       "Internal server error"
// @Router       /api/reservation/my-reservation [get]
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



// DeleteReservation godoc
// @Summary      Delete a reservation by ID
// @Description  Delete a reservation based on the provided reservation ID
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "Reservation ID"
// @Success      200  {object}  map[string]string    "Reservation deleted successfully"
// @Failure      400  {object}  response.ErrorResponse "Invalid reservation ID"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/reservation/{id} [delete]
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

// CreateReservation godoc
// @Summary      Create a new reservation
// @Description  Create a new reservation with the provided details
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        input  body      dto.Reservation          true  "Reservation creation details"
// @Success      201    {object}  response.ReservationResponse "Reservation created successfully"
// @Failure      400    {object}  response.ErrorResponse       "Invalid request body or validation failed"
// @Failure      409    {object}  response.ErrorResponse       "Table already reserved"
// @Failure      500    {object}  response.ErrorResponse       "Internal server error"
// @Router       /api/reservation [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	// var reservation models.Reservation
	var reservationDTO dto.Reservation

    // Bind dan validasi request body
    if err := c.ShouldBindJSON(&reservationDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
        return
    }

    // Validasi format waktu
    if _, err := time.Parse(time.RFC3339, reservationDTO.ReservationDateTime.Format(time.RFC3339)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation date format. Expected RFC3339 format."})
        return
    }

    // Validasi request menggunakan validator
    if err := h.validate.Struct(reservationDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

    // Convert DTO ke model
    reservationModel := models.Reservation{
        UserID:             reservationDTO.UserID,
        TableID:            reservationDTO.TableID,
        ReservationDateTime: reservationDTO.ReservationDateTime,
        NumberOfPeople:     reservationDTO.NumberOfPeople,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

	//Get email userLogin
	emailInterface, exists :=  c.Get("email")	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error" : "email tidak ditemukan di claims",
		})
		return
	}

	// Parse claims email to string
	userEmail, ok := emailInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Gagal konversi userEmail ke string",
		})
		return
	}


    // Panggil service untuk membuat reservasi
    if err := h.ReservationService.CreateReservation(&reservationModel, userEmail ); err != nil {
        switch {
        case strings.Contains(err.Error(), "already reserved"):
            c.JSON(http.StatusBadRequest, gin.H{"error": "Table is already reserved"})
        case strings.Contains(err.Error(), "failed to fetch table"):
            c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
        default:
            log.Printf("Error creating reservation: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
        }
        return
    }

    // Ambil detail reservasi yang baru dibuat
    createReservation, err := h.ReservationService.GetReservationDetail(reservationModel.ID)
    if err != nil {
        log.Printf("Error fetching reservation detail: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservation detail"})
        return
    }

    // Format response
    resp := response.ReservationResponse{
        ID: createReservation.ID,
        User: response.UserPreloadResponse{
            ID:    createReservation.UserID,
            Name:  createReservation.User.Name,
            Email: createReservation.User.Email,
        },
        Table: response.TableResponse{
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

    // Kirim response sukses
    c.JSON(http.StatusCreated, gin.H{
        "message": "Reservation created successfully",
        "data":    resp,
    })
}


// UpdateReservation godoc
// @Summary      Update a reservation by ID
// @Description  Update a reservation's details based on its ID
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id     path      int                      true  "Reservation ID"
// @Param        input  body      dto.UpdateReservation    true  "Updated reservation details"
// @Success      200    {object}  map[string]string        "Reservation updated successfully"
// @Failure      400    {object}  response.ErrorResponse   "Invalid request body or validation failed"
// @Failure      409    {object}  response.ErrorResponse   "Table and reservation time already used by another reservation"
// @Failure      500    {object}  response.ErrorResponse   "Internal server error"
// @Router       /api/reservation/{id} [put]
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