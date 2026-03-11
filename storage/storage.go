package storage

import "booking-be/models"

// Store defines the interface for data storage (DynamoDB-backed)
type Store interface {
	GetRooms() []models.Room
	GetRoomByID(id string) (*models.Room, bool)
	GetBookings() []models.Booking
	GetBookingByID(id string) (*models.Booking, bool)
	CreateBooking(booking models.Booking) models.Booking
}
