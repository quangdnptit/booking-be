package service

import (
	"context"

	"booking-be/models"
	"booking-be/repo"
)

type SeatService struct {
	seatRepo repo.SeatRepo
	priceCfg SeatPriceConfig
}

// NewSeatService creates a SeatService with optional seat price config (base + per-type modifier).
func NewSeatService(seatRepo repo.SeatRepo, priceCfg SeatPriceConfig) *SeatService {
	return &SeatService{
		seatRepo: seatRepo,
		priceCfg: priceCfg,
	}
}

// GenerateSeats computes each seat's price as base_price + modifier(seat_type), then batch-persists.
func (s *SeatService) GenerateSeats(ctx context.Context, basePrice float64, seats []models.Seat) error {
	for i := range seats {
		seatType := seats[i].SeatType
		if seatType == "" {
			seatType = models.SeatTypeStandard
		}
		seats[i].Price = s.priceCfg.SeatPrice(basePrice, seatType)
	}
	return s.seatRepo.GenerateSeats(ctx, seats)
}

// GetSeats returns all seats for a showtimeId
func (s *SeatService) GetSeats(ctx context.Context, showtimeId string) ([]models.Seat, error) {
	return s.seatRepo.GetByShowtimeID(ctx, showtimeId)
}
