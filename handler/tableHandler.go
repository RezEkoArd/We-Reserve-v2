package handler

import (
	"net/http"
	"strconv"
	"wereserve/dto"
	"wereserve/handler/response"
	"wereserve/models"
	"wereserve/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TableHandler struct {
	TableService *services.TableService
}

func NewTableHandler(tableService *services.TableService) *TableHandler {
	return &TableHandler{TableService: tableService}
}

func (h *TableHandler) GetListTable(c *gin.Context) {
	tables, err := h.TableService.GetAllTable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error()})
			return
	}

	// Convert ke response list
	var tableResponse []response.TableResponse
	for _, table := range tables{
		tableResponse = append(tableResponse, response.TableResponse{
			ID:        table.ID,
			TableName: table.TableName,
			Capacity:  table.Capacity,
			Status:    table.Status,
			CreatedAt: table.CreatedAt,
			UpdatedAt: table.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data" : tableResponse,
	})
}


func (h *TableHandler) GetTableByID(c *gin.Context) {
	//get Param
	ParamId := c.Param("id")
	id, err := strconv.Atoi(ParamId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Id tidak ditemukan"})
		return
	}

	//panggil service
	table, err := h.TableService.GetTableById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Data not Found"})
		return
	}

	resp := response.TableResponse{
		ID:        table.ID,
		TableName: table.TableName,
		Capacity:  table.Capacity,
		Status:    table.Status,
		CreatedAt: table.CreatedAt,
		UpdatedAt: table.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "data get successfully",
		"data" : resp,
	})
}

func (h *TableHandler) CreateTable(c *gin.Context) {
	var req dto.CreateTableRequest

	// bind json 
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}


	// Validation Json Req
	if err := dto.Validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Validator failed", "details" : errors})
		return
	}

	//Mapping DTO ke model
	table := models.Table{
		TableName: req.TableName,
		Capacity: req.Capacity,
		Status: req.Status,
	}

	// panggil service
	if err := h.TableService.CreateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	//return
	c.JSON(http.StatusCreated, gin.H{
		"message" : "Table created Has been Successfully",
		"data" : table,
	})
}

// Update Table

func (h *TableHandler) UpdateTable(c *gin.Context) {
	// Ambil ID dari paramater
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid ID"})
	}


	// Bind JSON request 
	var req dto.UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}


	// validation
	if err := dto.Validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _,err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}

		c.JSON(http.StatusBadRequest, gin.H{"error" : "Validation Error", "details" : errors})
		return
	}

	//Mapping DTO ke model
	table := models.Table{
		TableName: req.TableName,
		Capacity: req.Capacity,
		Status: req.Status,
	}

	if err := h.TableService.UpdateTable(id, table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Table update successfully"})
}

// delete Table
func (h *TableHandler) DeleteTable(c *gin.Context) {
	// Ambil id 
	ParamID := c.Param("id")
	id, err := strconv.Atoi(ParamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Id Tidak ditemukan"})
	}

	// bind JSON Request 
	err = h.TableService.DeleteTable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message" : "Table deleted successfully"})
}

func (h *TableHandler) GetTableByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Status is required"})
		return
	}

	// panggil service
	tables, err := h.TableService.GetTableByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors" : err.Error()})
		return
	}

	// Convert ke response list
	var tableResponse []response.TableResponse
	for _, table := range tables {
		tableResponse = append(tableResponse, response.TableResponse{
			ID:        table.ID,
			TableName: table.TableName,
			Capacity:  table.Capacity,
			Status:    table.Status,
			CreatedAt: table.CreatedAt,
			UpdatedAt: table.UpdatedAt,
		})
	}

	//response
	c.JSON(http.StatusOK, gin.H{"data" : tableResponse})
}
 