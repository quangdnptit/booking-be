package repo

import (
	"context"
	"errors"
	"fmt"

	"booking-be/models"
	"booking-be/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// PostgresProgramRepo loads movies and showtimes from the same Postgres DB as cms-booking (Flyway V1 + V2).
type PostgresProgramRepo struct {
	pool *storage.PostgresPool
}

func NewPostgresProgramRepo(pool *storage.PostgresPool) *PostgresProgramRepo {
	return &PostgresProgramRepo{pool: pool}
}

var ErrNotFound = errors.New("not found")

func scanMovieResponse(row pgx.Row) (models.MovieResponse, error) {
	var m models.MovieResponse
	var desc, genre, age, poster pgtype.Text
	err := row.Scan(
		&m.ID, &m.Title, &desc, &m.DurationMinutes, &genre, &age, &poster,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return m, err
	}
	if desc.Valid {
		m.Description = desc.String
	}
	if genre.Valid {
		m.Genre = genre.String
	}
	if age.Valid {
		m.AgeRating = age.String
	}
	if poster.Valid {
		m.PosterURL = poster.String
	}
	return m, nil
}

func (r *PostgresProgramRepo) ListMovies(ctx context.Context) ([]models.MovieResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, duration_minutes, genre, age_rating, poster_url, created_at, updated_at
		FROM movies
		ORDER BY title ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list movies: %w", err)
	}
	defer rows.Close()
	var out []models.MovieResponse
	for rows.Next() {
		m, err := scanMovieResponse(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *PostgresProgramRepo) GetMovieByID(ctx context.Context, id string) (*models.MovieResponse, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, title, description, duration_minutes, genre, age_rating, poster_url, created_at, updated_at
		FROM movies WHERE id = $1::uuid
	`, id)
	m, err := scanMovieResponse(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get movie: %w", err)
	}
	return &m, nil
}

func scanShowtimeResponse(row pgx.Row) (models.ShowtimeResponse, error) {
	var s models.ShowtimeResponse
	var pubAt pgtype.Timestamptz
	var price pgtype.Numeric
	err := row.Scan(
		&s.ID, &s.MovieID, &s.MovieTitle, &s.RoomID, &s.RoomName,
		&s.StartTime, &s.EndTime, &pubAt, &s.IsPublished, &price,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return s, err
	}
	if pubAt.Valid {
		t := pubAt.Time
		s.PublishedAt = &t
	}
	if price.Valid {
		f, _ := price.Float64Value()
		if f.Valid {
			s.BasePrice = f.Float64
		}
	}
	return s, nil
}

func (r *PostgresProgramRepo) ListShowtimes(ctx context.Context) ([]models.ShowtimeResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.movie_id, m.title, s.room_id, r.name,
		       s.start_time, s.end_time, s.published_at, s.is_published, s.base_price,
		       s.created_at, s.updated_at
		FROM showtimes s
		JOIN movies m ON m.id = s.movie_id
		JOIN rooms r ON r.id = s.room_id
		ORDER BY s.start_time ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list showtimes: %w", err)
	}
	defer rows.Close()
	var out []models.ShowtimeResponse
	for rows.Next() {
		s, err := scanShowtimeResponse(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PostgresProgramRepo) GetShowtimeByID(ctx context.Context, id string) (*models.ShowtimeResponse, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT s.id, s.movie_id, m.title, s.room_id, r.name,
		       s.start_time, s.end_time, s.published_at, s.is_published, s.base_price,
		       s.created_at, s.updated_at
		FROM showtimes s
		JOIN movies m ON m.id = s.movie_id
		JOIN rooms r ON r.id = s.room_id
		WHERE s.id = $1::uuid
	`, id)
	s, err := scanShowtimeResponse(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get showtime: %w", err)
	}
	return &s, nil
}

func scanTheaterResponse(row pgx.Row) (models.TheaterResponse, error) {
	var t models.TheaterResponse
	err := row.Scan(&t.ID, &t.Name, &t.Location, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

// ListTheaters mirrors TheaterService.findAll (theaters + V2 updated_at).
func (r *PostgresProgramRepo) ListTheaters(ctx context.Context) ([]models.TheaterResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, location, created_at, updated_at
		FROM theaters
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list theaters: %w", err)
	}
	defer rows.Close()
	var out []models.TheaterResponse
	for rows.Next() {
		t, err := scanTheaterResponse(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// GetTheaterByID mirrors TheaterService.findById.
func (r *PostgresProgramRepo) GetTheaterByID(ctx context.Context, id string) (*models.TheaterResponse, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, location, created_at, updated_at
		FROM theaters WHERE id = $1::uuid
	`, id)
	t, err := scanTheaterResponse(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get theater: %w", err)
	}
	return &t, nil
}

func scanRoomResponse(row pgx.Row) (models.RoomResponse, error) {
	var rm models.RoomResponse
	err := row.Scan(
		&rm.ID, &rm.TheaterID, &rm.TheaterName, &rm.Name,
		&rm.TotalSeats, &rm.TotalRows, &rm.SeatsPerRow,
		&rm.CreatedAt, &rm.UpdatedAt,
	)
	return rm, err
}

// ListRoomsByTheaterID mirrors RoomService.findByTheater.
func (r *PostgresProgramRepo) ListRoomsByTheaterID(ctx context.Context, theaterID string) ([]models.RoomResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT r.id, r.theater_id, t.name, r.name,
		       r.total_seats, r.total_rows, r.seats_per_row,
		       r.created_at, r.updated_at
		FROM rooms r
		JOIN theaters t ON t.id = r.theater_id
		WHERE r.theater_id = $1::uuid
		ORDER BY r.name ASC
	`, theaterID)
	if err != nil {
		return nil, fmt.Errorf("list rooms by theater: %w", err)
	}
	defer rows.Close()
	var out []models.RoomResponse
	for rows.Next() {
		rm, err := scanRoomResponse(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rm)
	}
	return out, rows.Err()
}
