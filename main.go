package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"booking-be/handlers"
	"booking-be/storage"
)

func main() {
	// Initialize storage
	store := storage.NewInMemoryStore()

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
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
