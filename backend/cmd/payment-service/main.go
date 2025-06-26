package main

import (
	"log"
	"net/http"
	"os"

	"concert-booking/internal/payment"
	"concert-booking/pkg/database"
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
	paymentRepo := payment.NewPaymentRepository(db)
	paymentService := payment.NewPaymentService(paymentRepo)
	paymentHandler := payment.NewPaymentHandler(paymentService)

	// Create router
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	// Routes
	api := router.Group("/api")
	{
		api.POST("/payments", paymentHandler.ProcessPayment)
	}

	// Start server
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Payment service running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
