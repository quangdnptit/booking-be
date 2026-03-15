package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/dynamo/v2"

	"booking-be/models"
	"booking-be/repomodel"
	"booking-be/view"
)

const (
	TableBookings    = "bookings"
	TableBookedSeats = "booked_seats"
)

// BalanceDeduction is applied in the same transaction as booking (deduct user balance).
type BalanceDeduction struct {
	UserID        string
	CurrentAmount float64
	NewAmount     float64
	UpdatedAt     string
}

func BookSeatsTransaction(
	ctx context.Context,
	db *dynamo.DB,
	booking models.Bookings,
	seats []models.Seat,
	deduct *BalanceDeduction,
) error {
	if len(seats) == 0 {
		return fmt.Errorf("no seats to book")
	}

	bookingTbl := db.Table(TableBookings)
	seatTbl := db.Table(TableBookedSeats)

	bookingRecord := view.BookingDomain2Repo(booking)

	tx := db.WriteTx()
	tx.Put(bookingTbl.Put(bookingRecord))

	now := time.Now().UTC().Format(time.RFC3339)
	available := string(models.SeatStatusAvailable)

	for i := range seats {
		oldUpdatedAt := seats[i].UpdatedAt
		seats[i].BookingID = booking.ID
		seats[i].SeatStatus = models.SeatStatusBooked
		seats[i].UpdatedAt = now
		rec := view.BookedSeatDomain2Repo(seats[i])
		tx.Put(seatTbl.Put(rec).If("'updated_at' = ? AND seat_status = ?", oldUpdatedAt, available))
	}

	if deduct != nil {
		usersTbl := db.Table(TableUsers)
		pk := repomodel.PKPrefixUser + deduct.UserID
		sk := repomodel.SKProfile
		tx.Update(usersTbl.Update("pk", pk).Range("sk", sk).
			If("'amount' = ?", deduct.CurrentAmount).
			Set("amount", deduct.NewAmount).
			Set("updated_at", deduct.UpdatedAt))
	}

	if err := tx.Run(ctx); err != nil {
		return fmt.Errorf("book seats transaction: %w", err)
	}
	return nil
}
