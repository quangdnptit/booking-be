package service

import (
	"context"

	"booking-be/models"
	"booking-be/repo"
)

type SeatService struct {
	seatRepo repo.SeatRepo
}

// NewSeatService creates a SeatService
func NewSeatService(seatRepo repo.SeatRepo) *SeatService {
	return &SeatService{
		seatRepo: seatRepo,
	}
}

// GenerateSeats batch-persists domain seats via the seat repo (convert + BatchWriteItem).
func (s *SeatService) GenerateSeats(ctx context.Context, seats []models.Seat) error {
	return s.seatRepo.GenerateSeats(ctx, seats)
}

// GetSeats returns all seats for a showtimeId
func (s *SeatService) GetSeats(ctx context.Context, showtimeId string) ([]models.Seat, error) {
	return s.seatRepo.GetByShowtimeID(ctx, showtimeId)
}
