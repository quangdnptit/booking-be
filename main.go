package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"booking-be/handlers"
	"booking-be/repo"
	"booking-be/service"
	"booking-be/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	ctx := context.Background()
	dynamoDBEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	// Single DynamoDB client
	log.Println("Initializing DynamoDB client...")
	dynamodbStore, err := storage.NewDynamoDBStore(ctx, dynamoDBEndpoint)
	if err != nil {
		log.Fatalf("failed to create DynamoDB client: %v", err)
	}

	// Repos
	bookingRepo := repo.NewDynamoBookingRepo(dynamodbStore, "bookings")
	bookedSeatRepo := repo.NewDynamoBookedSeatRepo(dynamodbStore, "booked_seats")

	// Services
	svc := service.NewService(bookingRepo, bookedSeatRepo)
	seatService := service.NewSeatService(bookedSeatRepo)

	// Handlers
	handler := handlers.NewHandler(svc)
	seatHandler := handlers.NewSeatHandler(seatService)

	// Setup router
	router := gin.Default()
	// Health endpoint
	router.GET("/api/v1/health", handler.HealthCheck)
	// Handle Seat endpoints
	router.POST("/api/v1/seats/generate-seats", seatHandler.GenerateSeats)

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
