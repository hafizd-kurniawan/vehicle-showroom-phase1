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

  // Initialize repositories
  userRepo := repository.NewUserRepository(db)
  sessionRepo := repository.NewSessionRepository(db)

  // Initialize use cases
  authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, cfg.JWT)

  // Initialize HTTP handlers
  authHandler := http.NewAuthHandler(authUsecase)

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
      auth.POST("/logout", authHandler.Logout)
      auth.GET("/me", authHandler.GetProfile)
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
