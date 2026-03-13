package models

import "github.com/google/uuid"

type SeatStatus string

const (
	SeatStatusUnknown     SeatStatus = "UNKNOWN"
	SeatStatusAvailable   SeatStatus = "AVAILABLE"
	SeatStatusUnAvailable SeatStatus = "UNAVAILABLE"
	SeatStatusBooked      SeatStatus = "BOOKED"
	SeatStatusLocked      SeatStatus = "LOCKED"
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
	ShowtimeID string     `json:"showtime_id"`
	SeatKey    string     `json:"seat_key"` // e.g. row#line
	RoomID     uuid.UUID  `json:"room_id"`
	SeatType   SeatType   `json:"seat_type"`
	BookingID  string     `json:"booking_id"`
	IsActive   string     `json:"is_active"`
	Price      float32    `json:"price"`
	SeatStatus SeatStatus `json:"seat_status"`
	CreatedAt  string     `json:"created_at,omitempty"`
	UpdatedAt  string     `json:"updated_at,omitempty"`
}
