package service

import (
	"booking-be/models"
	"booking-be/repo"
	"context"
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
	if err := s.seatRepo.GenerateSeats(ctx, seats); err != nil {
		return err
	}
	return nil
}
