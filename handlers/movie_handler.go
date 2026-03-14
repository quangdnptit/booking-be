package handlers

import (
	"errors"
	"net/http"

	"booking-be/internal/observability"
	"booking-be/models"
	"booking-be/repo"
	"booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProgramHandler serves GET /api/movies and GET /api/showtimes (same contract as cms-booking).
type ProgramHandler struct {
	svc *service.ProgramService
}

func NewProgramHandler(svc *service.ProgramService) *ProgramHandler {
	return &ProgramHandler{svc: svc}
}

func (h *ProgramHandler) ListMovies(c *gin.Context) {
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

func (h *ProgramHandler) GetMovieByID(c *gin.Context) {
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

func (h *ProgramHandler) ListShowtimes(c *gin.Context) {
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

func (h *ProgramHandler) GetShowtimeByID(c *gin.Context) {
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
