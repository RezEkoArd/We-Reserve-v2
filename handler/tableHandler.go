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

// GetListTable godoc
// @Summary      Get all tables
// @Description  Retrieve a list of all tables
// @Tags         tables
// @Accept       json
// @Produce      json
// @Success      200  {array}   response.TableResponse "List of tables retrieved successfully"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/tables [get]
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

// GetTableByID godoc
// @Summary      Get a table by ID
// @Description  Retrieve a table's information based on its ID
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "Table ID"
// @Success      200  {object}  response.TableResponse "Table retrieved successfully"
// @Failure      400  {object}  response.ErrorResponse "Invalid table ID"
// @Failure      404  {object}  response.ErrorResponse "Table not found"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/tables/{id} [get]
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


// CreateTable godoc
// @Summary      Create a new table
// @Description  Create a new table with the provided details
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateTableRequest  true  "Table creation details"
// @Success      201    {object}  map[string]interface{}  "Table created successfully"
// @Failure      400    {object}  response.ErrorResponse  "Invalid request body or validation failed"
// @Failure      500    {object}  response.ErrorResponse  "Internal server error"
// @Router       /api/tables [post]
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

// UpdateTable godoc
// @Summary      Update a table by ID
// @Description  Update a table's information based on its ID
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id     path      int                    true  "Table ID"
// @Param        input  body      dto.UpdateTableRequest  true  "Updated table details"
// @Success      200    {object}  map[string]string      "Table updated successfully"
// @Failure      400    {object}  response.ErrorResponse "Invalid request body or validation failed"
// @Failure      500    {object}  response.ErrorResponse "Internal server error"
// @Router       /api/tables/{id} [put]
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

// DeleteTable godoc
// @Summary      Delete a table by ID
// @Description  Delete a table based on the provided table ID
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Table ID"
// @Success      200  {object}  map[string]string "Table deleted successfully"
// @Failure      400  {object}  response.ErrorResponse "Invalid table ID"
// @Failure      500  {object}  response.ErrorResponse "Internal server error"
// @Router       /api/tables/{id} [delete]
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

// GetTableByStatus godoc
// @Summary      Get tables by status
// @Description  Retrieve a list of tables filtered by their status
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        status  query    string               true  "Table status (e.g., available, reserved, occupied)"
// @Success      200     {array}  response.TableResponse "List of tables retrieved successfully"
// @Failure      400     {object} response.ErrorResponse "Invalid status parameter"
// @Failure      500     {object} response.ErrorResponse "Internal server error"
// @Router       /api/tables/status [get]
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
 