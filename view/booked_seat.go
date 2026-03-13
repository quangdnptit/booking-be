package view

import (
	"booking-be/models"
	"booking-be/repomodel"

	"github.com/google/uuid"
)

func parseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// BookedSeatRepo2Domain maps a persistence record to the domain model
func BookedSeatRepo2Domain(r repomodel.BookedSeatRecord) models.Seat {
	statusStr := r.SeatStatus
	if statusStr == "" {
		statusStr = r.Status
	}
	seatStatus := models.SeatStatus(statusStr)
	if seatStatus == models.SeatStatusUnknown {
		seatStatus = models.SeatStatusAvailable
	}
	seatType := models.SeatType(r.SeatType)
	if seatType == models.SeatTypeUnknown {
		seatType = models.SeatTypeStandard
	}
	return models.Seat{
		ShowtimeID: r.ShowtimeID,
		SeatKey:    r.SeatKey,
		BookingID:  r.BookingID,
		RoomID:     parseUUID(r.RoomID),
		SeatType:   seatType,
		IsActive:   r.IsActive,
		Price:      r.Price,
		SeatStatus: seatStatus,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// BookedSeatDomain2Repo maps the domain model to a persistence record
func BookedSeatDomain2Repo(s models.Seat) repomodel.BookedSeatRecord {
	status := string(s.SeatStatus)
	if status == "" {
		status = string(models.SeatStatusAvailable)
	}
	seatType := string(s.SeatType)
	return repomodel.BookedSeatRecord{
		ShowtimeID: s.ShowtimeID,
		SeatKey:    s.SeatKey,
		BookingID:  s.BookingID,
		RoomID:     s.RoomID.String(),
		SeatType:   seatType,
		IsActive:   s.IsActive,
		Price:      s.Price,
		SeatStatus: status,
		Status:     status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
