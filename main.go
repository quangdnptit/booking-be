package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"booking-be/handlers"
	"booking-be/internal/observability"
	"booking-be/repo"
	"booking-be/service"
	"booking-be/storage"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Info().Str("event", "config_load").Msg(".env not found, using system env")
	}

	ctx := context.Background()
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	log.Info().Str("event", "dynamodb_init").Str("endpoint", endpoint).Msg("creating dynamo client")

	db, err := storage.NewDynamoDBStore(ctx, endpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("dynamodb client")
	}
	log.Info().Str("event", "dynamodb_ready").Msg("client ok")

	// Init dependencies
	bookingRepo := repo.NewDynamoBookingRepo(db)
	bookedSeatRepo := repo.NewDynamoBookedSeatRepo(db)
	svc := service.NewBookingService(bookingRepo, bookedSeatRepo, db)
	seatService := service.NewSeatService(bookedSeatRepo)
	handler := handlers.NewHandler(svc)
	seatHandler := handlers.NewSeatHandler(seatService)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Recovery())
	router.Use(observability.TracingMiddleware())

	router.GET("/api/v1/health", handler.HealthCheck)
	router.POST("/api/v1/seats/generate-seats", seatHandler.GenerateSeats)
	router.GET("/showtimes/:showtimeId/seats", seatHandler.GetSeats)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}
	log.Info().Str("event", "server_listen").Str("addr", port).Msg("listening")
	if err := router.Run(port); err != nil {
		log.Fatal().Err(err).Msg("server")
	}
}
