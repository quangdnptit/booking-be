package storage

import (
	"sync"

	"booking-be/models"
)

// Store defines the interface for data storage
type Store interface {
	GetRooms() []models.Room
	GetRoomByID(id string) (*models.Room, bool)
	GetBookings() []models.Booking
	GetBookingByID(id string) (*models.Booking, bool)
	CreateBooking(booking models.Booking) models.Booking
}

// InMemoryStore implements Store interface with in-memory storage
type InMemoryStore struct {
	mu       sync.RWMutex
	rooms    []models.Room
	bookings []models.Booking
}

// NewInMemoryStore creates and initializes a new in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		rooms: []models.Room{
			{ID: "room-1", Name: "Deluxe Suite", Capacity: 2, PricePerNight: 150.00},
			{ID: "room-2", Name: "Family Room", Capacity: 4, PricePerNight: 200.00},
		},
		bookings: []models.Booking{},
	}
}

// GetRooms returns all rooms
func (s *InMemoryStore) GetRooms() []models.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rooms
}

// GetRoomByID returns a room by its ID
func (s *InMemoryStore) GetRoomByID(id string) (*models.Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, room := range s.rooms {
		if room.ID == id {
			return &s.rooms[i], true
		}
	}
	return nil, false
}

// GetBookings returns all bookings
func (s *InMemoryStore) GetBookings() []models.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bookings
}

// GetBookingByID returns a booking by its ID
func (s *InMemoryStore) GetBookingByID(id string) (*models.Booking, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, booking := range s.bookings {
		if booking.ID == id {
			return &s.bookings[i], true
		}
	}
	return nil, false
}

// CreateBooking creates and stores a new booking
func (s *InMemoryStore) CreateBooking(booking models.Booking) models.Booking {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookings = append(s.bookings, booking)
	return booking
}
