package handler

import (
	"net/http"
	"strconv"
	"wereserve/dto"
	response "wereserve/handler/response"
	"wereserve/models"
	"wereserve/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	//Bind json request body ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid body Request"})
		return
	}

	// Validasi input menggunakan validator10
	if err := dto.Validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Validator failed", "details" : errors})
		return
	}

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
	}
	
	// Panggil service untuk register user
	if err := h.UserService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	
	//bind JSON Request body struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid Body Request"})
		return
	}


	// Validasi input using validator10
	if err := dto.Validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _,err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors" : "Validator failed", "details" : errors})
		return
	}

	//panggil Service untuk login user
	token, err := h.UserService.LoginUser(req.Email, req.Password)
	if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error" : err.Error()})
			return
	}

	//Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*3, "","", false, true)

	c.JSON(http.StatusOK, gin.H{"token" : token})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	// get id from Params
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "id tidak ditemukan" })
	}

	// Panggil service untuk menghapus user
	err = h.UserService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	users, err := h.UserService.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{"data" : users})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	// get id param
	ParamId := c.Param("id")
	id, err := strconv.Atoi(ParamId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Id tidak di temukan"})
		return
	}

	// panggil service
	user, err := h.UserService.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "Data not Found"})
		return
	}

	resp := response.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}

	// jika data ditemukan
	c.JSON(http.StatusOK, gin.H{
		"message" : "data get successfully",
		"data" : resp,
	})
}

// Update User
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse Id from url
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Id tidak di temukan"})
		return
	}

	//bind json body
	var user dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid request body"})
		return
	}

	if err := dto.Validate.Struct(user); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"error" : "Validator failed", "details" : errors})
		return
	}

	// Convert Dto to model
	reqBody := models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password: user.Password,
	}

	// panggil service
	err = h.UserService.UpdateUser(id, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	//response with success
	c.JSON(http.StatusOK, gin.H{"message" : "User Update Successfully"})

}