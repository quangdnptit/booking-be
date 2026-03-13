package repomodel

// BookingRecord is the persistence model for the DynamoDB bookings table.
// Keys: hash id, range created_at.
type BookingRecord struct {
	ID          string  `dynamo:"id"`
	UserID      string  `dynamo:"user_id"`
	ShowtimeID  string  `dynamo:"showtime_id"`
	TotalAmount float64 `dynamo:"total_amount"`
	Status      string  `dynamo:"status"`
	CreatedAt   string  `dynamo:"created_at"`
	UpdatedAt   string  `dynamo:"updated_at"`
}

// BookedSeatRecord is the persistence model for booked_seats.
// Table: pk showtime_id, sk seat_key. GSI: booking-seats-index on booking_id.
type BookedSeatRecord struct {
	ShowtimeID string  `dynamo:"showtime_id"`
	SeatKey    string  `dynamo:"seat_key"`
	BookingID  string  `dynamo:"booking_id"`
	RoomID     string  `dynamo:"room_id"`
	SeatType   string  `dynamo:"seat_type"`
	IsActive   string  `dynamo:"is_active"`
	Price      float32 `dynamo:"price"`
	SeatStatus string  `dynamo:"seat_status"`
	Status     string  `dynamo:"status"` // alias for Update expressions
	CreatedAt  string  `dynamo:"created_at"`
	UpdatedAt  string  `dynamo:"updated_at"`
}
