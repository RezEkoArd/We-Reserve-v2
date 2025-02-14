package handler

import (
	"net/http"
	"wereserve/dto"
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