package app

import (
	"wereserve/config"
	"wereserve/handler"
	"wereserve/middleware"
	"wereserve/repository"
	"wereserve/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal().Msgf("Error Connection to database: %v", err)
	}

	// inisialisasi handler routes dan service
	userRepo := repository.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)


	// Setup gin router
	r := gin.Default()

	// Public routes(tanpa middleware JWT)
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}
	
	//Endpoint yang memerlukan authentication dan role tertentu
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware()) // gunakan middleware jwt
	{
		// Contoh menggunakan jwt admin 
		api.DELETE("/users/:id", middleware.AdminOnlyMiddleware(), userHandler.DeleteUser)
		api.GET("/users", middleware.AdminOnlyMiddleware(), userHandler.GetAllUser)
		api.GET("/users/:id", middleware.AdminOnlyMiddleware(), userHandler.GetUserById)
		api.PUT("/users/:id", middleware.AdminOnlyMiddleware(), userHandler.UpdateUser)
	}

	



	r.Run(":8080")

}