package app

import (
	"net/http"
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

	// Protect routes (dengan menggunakan middleware JWT)
	// api := r.Group("/api")
	// api.Use(middleware.JWTAuthMiddleware())
	// {
	// 	api.GET("/profile", func ()  {
			
	// 	})
	// }



	//Endpoint yang memerlukan authentication dan role tertentu
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware()) // gunakan middleware jwt

	// Contoh menggunakan jwt admin 
	api.GET("/admin", middleware.RoleCheck("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Halo Admin!",
		})
	})


	r.Run(":8080")

}