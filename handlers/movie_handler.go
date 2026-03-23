package handlers

import (
	"errors"
	"net/http"

	"github.com/quangdnptit/booking-be/internal/observability"
	"github.com/quangdnptit/booking-be/models"
	"github.com/quangdnptit/booking-be/repo"
	"github.com/quangdnptit/booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// MovieHandler serves GET /api/movies, GET /api/showtimes
type MovieHandler struct {
	svc *service.MovieService
}

func NewMovieHandler(svc *service.MovieService) *MovieHandler {
	return &MovieHandler{svc: svc}
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	list, err := h.svc.ListMovies(c.Request.Context())
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "movies_list").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load movies"})
		return
	}
	if list == nil {
		list = []models.MovieResponse{}
	}
	log.Info().Str("trace_id", traceID).Str("event", "movies_list_ok").Int("n", len(list)).Send()
	c.JSON(http.StatusOK, list)
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	m, err := h.svc.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		log.Error().Str("trace_id", traceID).Str("event", "movie_get_by_id").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load movie"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *MovieHandler) ListShowtimes(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	list, err := h.svc.ListShowtimes(c.Request.Context())
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "showtimes_list").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load showtimes"})
		return
	}
	if list == nil {
		list = []models.ShowtimeResponse{}
	}
	log.Info().Str("trace_id", traceID).Str("event", "showtimes_list_ok").Int("n", len(list)).Send()
	c.JSON(http.StatusOK, list)
}

func (h *MovieHandler) GetShowtimeByID(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	s, err := h.svc.GetShowtimeByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		log.Error().Str("trace_id", traceID).Str("event", "showtime_get_by_id").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load showtime"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *MovieHandler) ListTheaters(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	list, err := h.svc.ListTheaters(c.Request.Context())
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "theaters_list").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load theaters"})
		return
	}
	if list == nil {
		list = []models.TheaterResponse{}
	}
	log.Info().Str("trace_id", traceID).Str("event", "theaters_list_ok").Int("n", len(list)).Send()
	c.JSON(http.StatusOK, list)
}

func (h *MovieHandler) GetTheaterByID(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	t, err := h.svc.GetTheaterByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		log.Error().Str("trace_id", traceID).Str("event", "theater_get_by_id").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load theater"})
		return
	}
	c.JSON(http.StatusOK, t)
}

// ListRoomsByTheater GET /api/rooms/theater/:theaterId — same as RoomController.findByTheater.
func (h *MovieHandler) ListRoomsByTheater(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	theaterID := c.Param("theaterId")
	if theaterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "theaterId required"})
		return
	}
	list, err := h.svc.ListRoomsByTheaterID(c.Request.Context(), theaterID)
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "rooms_by_theater").Str("theater_id", theaterID).Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load rooms"})
		return
	}
	if list == nil {
		list = []models.RoomResponse{}
	}
	log.Info().Str("trace_id", traceID).Str("event", "rooms_by_theater_ok").Str("theater_id", theaterID).Int("n", len(list)).Send()
	c.JSON(http.StatusOK, list)
}
