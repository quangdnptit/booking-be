package service

import (
	"os"
	"strconv"

	"github.com/quangdnptit/booking-be/models"
)

// SeatPriceConfig holds per-seat-type price modifiers (added to base price). Load from env.
type SeatPriceConfig struct {
	Standard   float64
	Premium    float64
	Wheelchair float64
	Unknown    float64
}

// LoadSeatPriceConfigFromEnv reads SEAT_PRICE_MODIFIER_* (defaults: Standard=0, Premium=5, Wheelchair=10).
func LoadSeatPriceConfigFromEnv() SeatPriceConfig {
	return SeatPriceConfig{
		Standard:   parseFloatEnv("SEAT_PRICE_MODIFIER_STANDARD", 0),
		Premium:    parseFloatEnv("SEAT_PRICE_MODIFIER_PREMIUM", 5),
		Wheelchair: parseFloatEnv("SEAT_PRICE_MODIFIER_WHEELCHAIR", 10),
		Unknown:    parseFloatEnv("SEAT_PRICE_MODIFIER_UNKNOWN", 0),
	}
}

func parseFloatEnv(key string, def float64) float64 {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}

// Modifier returns the additive modifier for the given seat type.
func (c SeatPriceConfig) Modifier(t models.SeatType) float64 {
	switch t {
	case models.SeatTypeStandard:
		return c.Standard
	case models.SeatTypePremium:
		return c.Premium
	case models.SeatTypeWheelChair:
		return c.Wheelchair
	default:
		return c.Unknown
	}
}

// SeatPrice returns base_price + modifier for the seat type.
func (c SeatPriceConfig) SeatPrice(basePrice float64, seatType models.SeatType) float32 {
	return float32(basePrice + c.Modifier(seatType))
}
