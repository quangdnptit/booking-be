package models

import "github.com/google/uuid"

type SeatStatus string

const (
	SeatStatusUnknown   SeatStatus = "UNKNOWN"
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusBooked    SeatStatus = "BOOKED"
	SeatStatusLocked    SeatStatus = "LOCKED"
)

type SeatType string

const (
	SeatTypeUnknown    SeatType = "UNKNOWN"
	SeatTypeStandard   SeatType = "STANDARD"
	SeatTypePremium    SeatType = "PREMIUM"
	SeatTypeWheelChair SeatType = "WHEELCHAIR"
)

// Bookings is the domain model for a showtime booking
type Bookings struct {
	UserID      string
	ShowtimeID  string
	TotalAmount float64
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

// Seat is the domain model for a booked seat
type Seat struct {
	ShowtimeID string
	SeatKey    string //{row#line}
	RoomID     uuid.UUID
	SeatType   SeatType
	BookingID  string
	IsActive   string
	Price      float32
	SeatStatus SeatStatus
	CreatedAt  string
	UpdatedAt  string
}
