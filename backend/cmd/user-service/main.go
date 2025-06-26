package main

import (
	"log"
	"net/http"
	"os"

	"concert-booking/pkg/database"
	"concert-booking/pkg/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto migrate models
	if err := database.Migrate(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Initialize services
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Create router
	router := gin.Default()

	// Routes
	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// Start server
	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("User service running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
