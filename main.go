package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"booking-be/handlers"
	"booking-be/repo"
	"booking-be/service"
	"booking-be/storage"
)

func main() {
	ctx := context.Background()
	dynamoDBEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	// Single DynamoDB client (DI)
	log.Println("Initializing DynamoDB client...")
	dynamodbStore, err := storage.NewDynamoDBStore(ctx, dynamoDBEndpoint)
	if err != nil {
		log.Fatalf("failed to create DynamoDB client: %v", err)
	}

	if err != nil {
		log.Fatalf("failed to initialize DynamoDB store: %v", err)
	}

	// Repos (showtime bookings + booked seats)
	bookingRepo := repo.NewDynamoBookingRepo(dynamodbStore, "bookings")
	bookedSeatRepo := repo.NewDynamoBookedSeatRepo(dynamodbStore, "booked_seats")

	// Service (DI: store + repos)
	svc := service.NewService(bookingRepo, bookedSeatRepo)

	// Handler (DI: service)
	handler := handlers.NewHandler(svc)
	// Setup router
	router := gin.Default()

	// Health endpoint
	router.GET("/health", handler.HealthCheck)
	router.GET("/health", handler.HealthCheck)

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
