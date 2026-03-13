package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	dynamo "github.com/guregu/dynamo/v2"

	"booking-be/models"
	"booking-be/repo"
)

const fixedSeatPrice = 20.0

// BookingService sits between handlers and persistence.
type BookingService struct {
	bookingRepo repo.BookingRepo
	seatRepo    repo.SeatRepo
	db          *dynamo.DB
}

// NewBookingService creates a service; db is used for atomic BookSeats (table names: repo.TableBookings / repo.TableBookedSeats).
func NewBookingService(bookingRepo repo.BookingRepo, seatRepo repo.SeatRepo, db *dynamo.DB) *BookingService {
	return &BookingService{bookingRepo: bookingRepo, seatRepo: seatRepo, db: db}
}

// BookSeats loads seats, checks availability, builds booking (fixed price per seat), then TransactWriteItems: booking + seat puts.
func (s *BookingService) BookSeats(ctx context.Context, req models.SeatsBookingRequest) error {
	if req.UserID == "" || req.ShowtimeID == "" {
		return fmt.Errorf("user_id and showtime_id are required")
	}
	if len(req.SeatKeys) == 0 {
		return fmt.Errorf("at least one seat_key is required")
	}
	if s.db == nil {
		return fmt.Errorf("dynamo db not configured for transactions")
	}

	seats, err := s.seatRepo.GetByShowtimeIDAndSeatKeys(ctx, req.ShowtimeID, req.SeatKeys)
	if err != nil {
		return fmt.Errorf("failed to load seats: %w", err)
	}
	keySeatMap := make(map[string]models.Seat, len(seats))
	for _, st := range seats {
		keySeatMap[st.SeatKey] = st
	}
	ordered, err := validateSeatsForBooking(req.ShowtimeID, req.SeatKeys, keySeatMap)
	if err != nil {
		return err
	}

	ts := time.Now().UTC().Format(time.RFC3339)
	bookingID := uuid.New().String()
	booking := models.Bookings{
		ID:          bookingID,
		UserID:      req.UserID,
		ShowtimeID:  req.ShowtimeID,
		TotalAmount: fixedSeatPrice * float64(len(ordered)),
		Status:      "CONFIRMED",
		CreatedAt:   ts,
		UpdatedAt:   ts,
	}

	if err := repo.BookSeatsTransaction(ctx, s.db, booking, ordered); err != nil {
		return err
	}
	return nil
}

func validateSeatsForBooking(showtimeID string, seatKeys []string, keySeatMap map[string]models.Seat) ([]models.Seat, error) {
	validSeat := make([]models.Seat, 0, len(seatKeys))
	for _, key := range seatKeys {
		seat, ok := keySeatMap[key]
		if !ok {
			return nil, fmt.Errorf("seat %q not found for showtime %q", key, showtimeID)
		}
		if seat.BookingID != "" {
			return nil, fmt.Errorf("seat %q already held or booked", key)
		}
		status := seat.SeatStatus
		if status == models.SeatStatusUnknown {
			status = models.SeatStatusAvailable
		}
		if status != models.SeatStatusAvailable {
			return nil, fmt.Errorf("seat %q is not available (status=%s)", key, seat.SeatStatus)
		}
		validSeat = append(validSeat, seat)
	}
	return validSeat, nil
}
