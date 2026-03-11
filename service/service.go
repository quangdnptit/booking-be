package service

import (
	"context"

	"booking-be/models"
	"booking-be/repo"
)

// Service sits between handlers and persistence (storage + repos).
// All dependencies are injected via the constructor.
type Service struct {
	bookingRepo    repo.BookingRepo
	bookedSeatRepo repo.BookedSeatRepo
}

// NewService creates a Service with the given dependencies (DI)
func NewService(bookingRepo repo.BookingRepo, bookedSeatRepo repo.BookedSeatRepo) *Service {
	return &Service{
		bookingRepo:    bookingRepo,
		bookedSeatRepo: bookedSeatRepo,
	}
}

// --- Showtime booking (via repo) ---
func (s *Service) GetShowtimeBookingByID(ctx context.Context, id string) (*models.Bookings, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func (s *Service) GetShowtimeBookingsByUserID(ctx context.Context, userID string) ([]models.Bookings, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}

func (s *Service) GetShowtimeBookingsByShowtimeID(ctx context.Context, showtimeID string) ([]models.Bookings, error) {
	return s.bookingRepo.GetByShowtimeID(ctx, showtimeID)
}

func (s *Service) CreateShowtimeBooking(ctx context.Context, booking models.Bookings) error {
	return s.bookingRepo.Create(ctx, booking)
}

func (s *Service) UpdateShowtimeBooking(ctx context.Context, booking models.Bookings) error {
	return s.bookingRepo.Update(ctx, booking)
}

func (s *Service) UpdateShowtimeBookingStatus(ctx context.Context, id, status string) error {
	return s.bookingRepo.UpdateStatus(ctx, id, status)
}

// --- Booked seats (via repo) ---

func (s *Service) GetBookedSeatByShowtimeIDAndSeatKey(ctx context.Context, showtimeID, seatKey string) (*models.BookedSeat, error) {
	return s.bookedSeatRepo.GetByShowtimeIDAndSeatKey(ctx, showtimeID, seatKey)
}

func (s *Service) GetBookedSeatsByBookingID(ctx context.Context, bookingID string) ([]models.BookedSeat, error) {
	return s.bookedSeatRepo.GetByBookingID(ctx, bookingID)
}

func (s *Service) GetBookedSeatsByShowtimeID(ctx context.Context, showtimeID string) ([]models.BookedSeat, error) {
	return s.bookedSeatRepo.GetByShowtimeID(ctx, showtimeID)
}

func (s *Service) CreateBookedSeat(ctx context.Context, seat models.BookedSeat) error {
	return s.bookedSeatRepo.Create(ctx, seat)
}

func (s *Service) UpdateBookedSeat(ctx context.Context, seat models.BookedSeat) error {
	return s.bookedSeatRepo.Update(ctx, seat)
}

func (s *Service) UpdateBookedSeatStatusByKey(ctx context.Context, showtimeID, seatKey, status string) error {
	return s.bookedSeatRepo.UpdateStatusByKey(ctx, showtimeID, seatKey, status)
}
