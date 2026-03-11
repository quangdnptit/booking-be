package models

// Room represents a hotel room
type Room struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Capacity      int     `json:"capacity"`
	PricePerNight float64 `json:"price_per_night"`
}

// Booking represents a room booking
type Booking struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

// CreateBookingRequest represents the request payload for creating a booking
type CreateBookingRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// Bookings is the domain model for a showtime booking
type Bookings struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	ShowtimeID string `json:"showtime_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// BookedSeat is the domain model for a booked seat
type BookedSeat struct {
	ID         string `json:"id"`
	BookingID  string `json:"booking_id"`
	ShowtimeID string `json:"showtime_id"`
	SeatID     string `json:"seat_id"`
	Status     string `json:"status"`
}
