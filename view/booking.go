package view

import (
	"booking-be/models"
	"booking-be/repomodel"
)

// BookingRepo2Domain maps a repo (persistence) record to the domain model
func BookingRepo2Domain(r repomodel.BookingRecord) models.Bookings {
	return models.Bookings{
		ID:         r.ID,
		UserID:     r.UserID,
		ShowtimeID: r.ShowtimeID,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
	}
}

// BookingDomain2Repo maps the domain model to a repo (persistence) record
func BookingDomain2Repo(b models.Bookings) repomodel.BookingRecord {
	return repomodel.BookingRecord{
		ID:         b.ID,
		UserID:     b.UserID,
		ShowtimeID: b.ShowtimeID,
		Status:     b.Status,
		CreatedAt:  b.CreatedAt,
	}
}
