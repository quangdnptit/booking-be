package service

import (
	"booking-be/models"
	"booking-be/repo"
	"context"
)

// Service sits between handlers and persistence
// All dependencies are injected via the constructor.
type BookingService struct {
	bookingRepo repo.BookingRepo
	seatRepo    repo.SeatRepo
}

// NewBookingService creates a Service with the given dependencies (DI)
func NewBookingService(bookingRepo repo.BookingRepo, seatRepo repo.SeatRepo) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		seatRepo:    seatRepo,
	}
}

// BookSeats core function of seats booking
func (s *BookingService) BookSeats(ctx context.Context, req models.SeatsBookingRequest) error {
	seatKeys := req.SeatKeys
	showTimeId := req.ShowtimeID

	_, err := s.seatRepo.GetByShowtimeIDAndSeatKeys(ctx, showTimeId, seatKeys)
	if err != nil {
		return err
	}

	return nil
}
