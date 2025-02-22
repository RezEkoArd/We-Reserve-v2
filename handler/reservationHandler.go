package handler

import (
	"net/http"
	"strconv"
	"wereserve/dto"
	"wereserve/handler/response"
	"wereserve/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReservationHandler struct {
	ReservationService *services.ReservationService
}

func NewReservationsHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{ReservationService: reservationService}
}


func (h *ReservationHandler) GetListReservation(c *gin.Context) {
	// fetchAll reservation from the service
	reservations, err := h.ReservationService.GetAllReservation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	//Convert reservation to response format
	var reservationResp []response.ReservationResponse
	for _,resp := range reservations {
		reservationResp = append(reservationResp, response.ReservationResponse{
			ID:             resp.ID,
			User:           response.UserPreloadResponse{
				ID:        resp.User.ID,
				Name:      resp.User.Name,
				Email:     resp.User.Email,
			},
			Table:          response.TableResponse{
				ID:        resp.Table.ID,
				TableName: resp.Table.TableName,
				Capacity:  resp.Table.Capacity,
				Status:    resp.Table.Status,
			},
			Date:           resp.Date,
			Time:           resp.Time,
			NumberOfPeople: resp.NumberOfPeople,
			CreatedAt:      resp.CreatedAt,
			UpdatedAt:      resp.UpdatedAt,
		})
	}

	// return response
	c.JSON(http.StatusOK, reservationResp)
}

func (h *ReservationHandler) GetReservationByUser(c *gin.Context) {
	// get userID
	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "User not authenticate"})
		return
	}
	
	//pastikan email itu string
	userEmailStr, ok := userEmail.(string)
	if !ok {
		c.JSON(http.StatusBadRequest,gin.H{"error" : "Invalid email user"})
	}

	// panggil service 
	reservation, err  := h.ReservationService.GetReservationByUsers(userEmailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	// Konversi ke response format
    var reservationResp []response.ReservationResponse
    for _, resp := range reservation {
        reservationResp = append(reservationResp, response.ReservationResponse{
            ID:             resp.ID,
            User: response.UserPreloadResponse{
                ID:    resp.User.ID,
                Name:  resp.User.Name,
                Email: resp.User.Email,
            },
			Table: response.TableResponse{
				ID:        resp.Table.ID,
				TableName: resp.Table.TableName,
				Capacity:  resp.Table.Capacity,
				Status:    resp.Table.Status,
			},
            Date:           resp.Date,
            Time:           resp.Time,
            NumberOfPeople: resp.NumberOfPeople,
            CreatedAt:      resp.CreatedAt,
            UpdatedAt:      resp.UpdatedAt,
        })
    }

    // Return the response
    c.JSON(http.StatusOK, reservationResp)
}

func (h *ReservationHandler) GetReservationById(c *gin.Context) {
	// get id
	ParamId := c.Param("id")
	id, err := strconv.Atoi(ParamId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid ID"})
	}

	// panggil service
	reservation, err := h.ReservationService.GetReservationDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Data Not Found"})
		return
	}

	resp := response.ReservationResponse{
		ID:             reservation.ID,
		User:           response.UserPreloadResponse{
			ID:        reservation.User.ID,
			Name:      reservation.User.Name,
			Email:     reservation.User.Email,
		},
		Table:          response.TableResponse{
			ID:        reservation.Table.ID,
			TableName: reservation.Table.TableName,
			Capacity:  reservation.Table.Capacity,
			Status:    reservation.Table.Status,
		},
		Date:           reservation.Date,
		Time:           reservation.Time,
		NumberOfPeople: reservation.NumberOfPeople,
		CreatedAt:      reservation.CreatedAt,
		UpdatedAt:      reservation.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "data get successfully",
		"data" : resp,
	})
}

func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid Id Reservation"})
	}

	// service panggil
	err = h.ReservationService.DeleteReservation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error" : "Reservation deleted successfully"})
}

func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	// Ambil ID reservasi
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "reservation id tidak ditemukan"})
		return
	}

	//bind req JSON
	var req dto.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}


	// Panggil service untuk update reservasi
    if err := h.ReservationService.UpdateReservation(id, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Berikan respons sukses
    c.JSON(http.StatusOK, gin.H{"message": "reservation updated successfully"})

}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req dto.CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	// validation 
	if err := dto.Validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Validator failed", 
			"details" : errors,
		})
		return
	}

	// panggil service
	if err := h.ReservationService.CreateReservation(
		req.UserID,
        req.TableID,
        req.Date,
        req.Time,
        req.NumberOfPeople,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message" : "Reservation Created has been Successfully",
		"data" : req,
	})
}