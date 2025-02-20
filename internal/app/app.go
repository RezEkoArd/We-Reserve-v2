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

	// inisialisasi handler routes dan service user
	userRepo := repository.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)


	// initilitaions Table 
	tableRepo := repository.NewTableRepository(db.DB)
	tableService := services.NewTableService(tableRepo)
	tableHandler := handler.NewTableHandler(tableService)


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
		api.DELETE("/users/:id", middleware.RoleCheck("admin"), userHandler.DeleteUser)
		api.GET("/users", middleware.RoleCheck("customer","admin"), userHandler.GetAllUser)
		api.GET("/users/:id", middleware.RoleCheck("customer","admin"), userHandler.GetUserById)
		api.PUT("/users/:id", middleware.RoleCheck("customer","admin"), userHandler.UpdateUser)

		api.GET("/tables", middleware.RoleCheck("customer", "admin"), tableHandler.GetListTable)
		api.GET("/tables/:id", middleware.RoleCheck("customer", "admin"),tableHandler.GetTableByID)
		api.GET("/tables/status", middleware.RoleCheck("customer", "admin"),tableHandler.GetTableByStatus)
		api.POST("/tables", middleware.RoleCheck("admin"),tableHandler.CreateTable)
		api.PUT("/tables/:id", middleware.RoleCheck("admin"),tableHandler.UpdateTable)
		api.DELETE("/tables/:id", middleware.RoleCheck("admin"),tableHandler.DeleteTable)
	}

	r.Run(":8080")

}