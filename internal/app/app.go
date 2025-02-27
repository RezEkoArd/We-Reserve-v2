package app

import (
	"time"
	"wereserve/config"
	"wereserve/handler"
	"wereserve/middleware"
	"wereserve/repository"
	"wereserve/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	// Swagger
	_ "wereserve/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			We Reserve
//	@version		1.0
//	@description	API untuk mengelola reservasi meja makan yang berada di restaurant
//	@description	Author: Rezekoard
//	@contact.name	RezkyEkoArd.
//	@contact.url	https://github.com/rezekoard

//	@host
//	@BasePath					/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@security					BearerAuth


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

	// inisilisasi Reservation
	reservationRepo := repository.NewReservationRepository(db.DB)
	reservationService := services.NewReservationService(reservationRepo, tableRepo)
	reservationHandler := handler.NewReservationsHandler(reservationService)


	// Setup gin router
	r := gin.Default()

	//cors Config
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// route Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes(tanpa middleware JWT)
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
		public.GET("/users", userHandler.GetAllUser)
		public.GET("/tables",tableHandler.GetListTable)

	}
	
	//Endpoint yang memerlukan authentication dan role tertentu
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware()) // gunakan middleware jwt
	{
		// Contoh menggunakan jwt admin 
		api.DELETE("/users/:id", middleware.RoleCheck("admin"), userHandler.DeleteUser)
		api.GET("/users/:id", middleware.RoleCheck("customer","admin"), userHandler.GetUserById)
		api.PUT("/users/:id", middleware.RoleCheck("customer","admin"), userHandler.UpdateUser)

		api.GET("/tables", middleware.RoleCheck("customer", "admin"), tableHandler.GetListTable)
		api.GET("/tables/:id", middleware.RoleCheck("customer", "admin"),tableHandler.GetTableByID)
		api.GET("/tables/status", middleware.RoleCheck("customer", "admin"),tableHandler.GetTableByStatus)
		api.POST("/tables", middleware.RoleCheck("admin"),tableHandler.CreateTable)
		api.PUT("/tables/:id", middleware.RoleCheck("admin"),tableHandler.UpdateTable)
		api.DELETE("/tables/:id", middleware.RoleCheck("admin"),tableHandler.DeleteTable)

		api.GET("/reservation", middleware.RoleCheck("admin"), reservationHandler.GetAllReservation)
		api.GET("/reservation/:id", middleware.RoleCheck("customer", "admin"),reservationHandler.GetReservationDetail)
		api.GET("/reservation/my-reservation", middleware.RoleCheck("customer", "admin"),reservationHandler.GetReservationByUserLogin)
		api.POST("/reservation", middleware.RoleCheck("customer","admin"),reservationHandler.CreateReservation)
		api.PUT("/reservation/:id", middleware.RoleCheck("admin"),reservationHandler.UpdateReservation)
		api.DELETE("/reservation/:id", middleware.RoleCheck("admin"),reservationHandler.DeleteReservation)
	}
	r.Run(":8080")

}