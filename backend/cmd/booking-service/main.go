package main

import (
	"log"
	"net/http"
	"os"

	"concert-booking/pkg/database"
	"concert-booking/pkg/internal/booking"
	"concert-booking/pkg/middleware"

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

	// Initialize services
	bookingRepo := booking.NewBookingRepository(db)
	bookingService := booking.NewBookingService(bookingRepo)
	bookingHandler := booking.NewBookingHandler(bookingService)

	// Create router
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	// Routes
	api := router.Group("/api")
	{
		api.POST("/bookings", bookingHandler.CreateBooking)
	}

	// Start server
	port := os.Getenv("BOOKING_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Booking service running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
