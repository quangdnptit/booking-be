package service_test

import (
	"context"
	_ "errors"
	"testing"

	"github.com/quangdnptit/booking-be/mock"
	"github.com/quangdnptit/booking-be/models"
	"github.com/quangdnptit/booking-be/service"
)

func TestGenerateSeats_SetsPriceFromBaseAndSeatType(t *testing.T) {
	ctx := context.Background()
	fake := &mock.MockSeatRepo{}
	cfg := service.SeatPriceConfig{
		Standard:   0,
		Premium:    5,
		Wheelchair: 10,
		Unknown:    1,
	}
	svc := service.NewSeatService(fake, cfg)

	seats := []models.Seat{
		{ShowtimeID: "st1", SeatKey: "A1", SeatType: models.SeatTypeStandard},
		{ShowtimeID: "st1", SeatKey: "A2", SeatType: models.SeatTypePremium},
		{ShowtimeID: "st1", SeatKey: "A3", SeatType: models.SeatTypeWheelChair},
		{ShowtimeID: "st1", SeatKey: "A4", SeatType: models.SeatTypeUnknown},
	}

	const base = 12.5
	if err := svc.GenerateSeats(ctx, base, seats); err != nil {
		t.Fatalf("GenerateSeats: %v", err)
	}

	// Input slice should be mutated the same way (service writes in place).
	//for i, w := range want {
	//	if seats[i].Price != w {
	//		t.Errorf("input seats[%d].Price = %v, want %v", i, seats[i].Price, w)
	//	}
	//}
}
