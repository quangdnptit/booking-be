package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"booking-be/handlers"
	"booking-be/storage"
)

func main() {
	dynamoDBEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	log.Println("Initializing DynamoDB store...")
	store, err := storage.NewDynamoDBStore(context.Background(), dynamoDBEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize DynamoDB store: %v", err)
	}
	log.Println("DynamoDB store initialized successfully")

	// Initialize handlers
	handler := handlers.NewHandler(store)

	// Setup router
	router := gin.Default()

	// Health endpoint
	router.GET("/health", handler.HealthCheck)

	// Room endpoints
	router.GET("/rooms", handler.ListRooms)
	router.GET("/rooms/:id", handler.GetRoom)
	router.GET("/rooms/:id/availability", handler.CheckAvailability)

	// Booking endpoints
	router.POST("/bookings", handler.CreateBooking)
	router.GET("/bookings", handler.ListBookings)
	router.GET("/bookings/:id", handler.GetBooking)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Starting server on %s...", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
