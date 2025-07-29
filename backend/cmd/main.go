package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"vehicle-showroom/internal/config"
	"vehicle-showroom/internal/database"
	"vehicle-showroom/internal/delivery/http"
	"vehicle-showroom/internal/repository"
	"vehicle-showroom/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize config
	cfg := config.New()

	// Initialize database
	db, err := database.NewPostgreSQL(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed demo data
	if err := database.SeedDemoData(db); err != nil {
		log.Fatal("Failed to seed demo data:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	reportRepo := repository.NewReportRepository(db)
	sparePartRepo := repository.NewSparePartRepository(db)
	repairRepo := repository.NewRepairRepository(db)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, cfg.JWT)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	vehicleUsecase := usecase.NewVehicleUsecase(vehicleRepo, customerRepo)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, vehicleRepo, customerRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo)
	sparePartUsecase := usecase.NewSparePartUsecase(sparePartRepo)
	repairUsecase := usecase.NewRepairUsecase(repairRepo, vehicleRepo, sparePartRepo)

	// Initialize HTTP handlers
	authHandler := http.NewAuthHandler(authUsecase)
	customerHandler := http.NewCustomerHandler(customerUsecase)
	vehicleHandler := http.NewVehicleHandler(vehicleUsecase)
	transactionHandler := http.NewTransactionHandler(transactionUsecase)
	reportHandler := http.NewReportHandler(reportUsecase)
	sparePartHandler := http.NewSparePartHandler(sparePartUsecase)
	repairHandler := http.NewRepairHandler(repairUsecase)

	// Initialize Middleware
	authMiddleware := http.AuthMiddleware(authUsecase)

	// Setup router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authMiddleware, authHandler.Logout)
			auth.GET("/me", authMiddleware, authHandler.GetProfile)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			customers := protected.Group("/customers")
			{
				customers.GET("", customerHandler.List)
				customers.POST("", http.RoleMiddleware("admin", "cashier"), customerHandler.Create)
				customers.GET("/:id", customerHandler.GetByID)
				customers.PUT("/:id", http.RoleMiddleware("admin", "cashier"), customerHandler.Update)
				customers.DELETE("/:id", http.RoleMiddleware("admin"), customerHandler.Delete)
			}

			vehicles := protected.Group("/vehicles")
			{
				vehicles.GET("", vehicleHandler.List)
				vehicles.POST("", http.RoleMiddleware("admin", "cashier"), vehicleHandler.Create)
				vehicles.GET("/:id", vehicleHandler.GetByID)
				vehicles.PUT("/:id", http.RoleMiddleware("admin", "cashier", "mechanic"), vehicleHandler.Update)
				vehicles.DELETE("/:id", http.RoleMiddleware("admin"), vehicleHandler.Delete)
				vehicles.PUT("/:id/status", http.RoleMiddleware("admin", "cashier", "mechanic"), vehicleHandler.UpdateStatus)
			}

			transactions := protected.Group("/transactions")
			{
				purchases := transactions.Group("/purchases")
				purchases.Use(http.RoleMiddleware("admin", "cashier"))
				{
					purchases.GET("", transactionHandler.ListPurchases)
					purchases.POST("", transactionHandler.CreatePurchase)
					purchases.GET("/:id", transactionHandler.GetPurchaseByID)
				}

				sales := transactions.Group("/sales")
				sales.Use(http.RoleMiddleware("admin", "cashier"))
				{
					sales.GET("", transactionHandler.ListSales)
					sales.POST("", transactionHandler.CreateSales)
					sales.GET("/:id", transactionHandler.GetSalesByID)
				}
			}

			dashboard := protected.Group("/dashboard")
			dashboard.Use(http.RoleMiddleware("admin"))
			{
				dashboard.GET("/stats", transactionHandler.GetDashboardStats)
			}

			reports := protected.Group("/reports")
			reports.Use(http.RoleMiddleware("admin"))
			{
				reports.GET("/profitability", reportHandler.GetVehicleProfitability)
				reports.GET("/sales", reportHandler.GetSalesReport)
				reports.GET("/purchases", reportHandler.GetPurchaseReport)
			}

			spareParts := protected.Group("/spare-parts")
			spareParts.Use(http.RoleMiddleware("admin", "mechanic"))
			{
				spareParts.GET("", sparePartHandler.List)
				spareParts.POST("", sparePartHandler.Create)
				spareParts.GET("/:id", sparePartHandler.GetByID)
				spareParts.PUT("/:id", sparePartHandler.Update)
				spareParts.DELETE("/:id", http.RoleMiddleware("admin"), sparePartHandler.Delete)
			}

			repairs := protected.Group("/repairs")
			repairs.Use(http.RoleMiddleware("admin", "mechanic"))
			{
				repairs.GET("", repairHandler.List)
				repairs.POST("", repairHandler.Create)
				repairs.GET("/:id", repairHandler.GetByID)
				repairs.PUT("/:id", repairHandler.Update)
				repairs.PUT("/:id/status", repairHandler.UpdateStatus)
				repairs.POST("/:id/parts", repairHandler.AddPart)
				repairs.DELETE("/:id/parts/:partId", repairHandler.RemovePart)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
