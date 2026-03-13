package repomodel

// BookingRecord is the persistence model for the DynamoDB bookings table.
type BookingRecord struct {
	ID          string  `dynamodbav:"id,omitempty"`
	UserID      string  `dynamodbav:"user_id"`
	ShowtimeID  string  `dynamodbav:"showtime_id"`
	TotalAmount float64 `dynamodbav:"total_amount,omitempty"`
	Status      string  `dynamodbav:"status"`
	CreatedAt   string  `dynamodbav:"created_at"`
	UpdatedAt   string  `dynamodbav:"updated_at,omitempty"`
}

// BookedSeatRecord is the persistence model for the DynamoDB booked_seats table.
// Table key: pk=showtime_id, sk=seat_key.
// seat_type: models.SeatType (STANDARD | PREMIUM | WHEELCHAIR) as string.
// seat_status: models.SeatStatus (AVAILABLE | BOOKED | LOCKED) as string.
type BookedSeatRecord struct {
	ShowtimeID string  `dynamodbav:"showtime_id"`
	SeatKey    string  `dynamodbav:"seat_key"`
	BookingID  string  `dynamodbav:"booking_id,omitempty"`
	RoomID     string  `dynamodbav:"room_id,omitempty"`
	SeatType   string  `dynamodbav:"seat_type,omitempty"` // SeatType enum
	IsActive   string  `dynamodbav:"is_active,omitempty"`
	Price      float32 `dynamodbav:"price,omitempty"`
	SeatStatus string  `dynamodbav:"seat_status,omitempty"` // SeatStatus enum
	Status     string  `dynamodbav:"status,omitempty"`      // mirror for UpdateItem expressions using #status
	CreatedAt  string  `dynamodbav:"created_at,omitempty"`
	UpdatedAt  string  `dynamodbav:"updated_at,omitempty"`
}
