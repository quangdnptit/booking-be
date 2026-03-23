package main

import (
	"context"
	"github.com/quangdnptit/booking-be/internal/observability"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/quangdnptit/booking-be/handlers"
	"github.com/quangdnptit/booking-be/internal/auth"
	"github.com/quangdnptit/booking-be/repo"
	"github.com/quangdnptit/booking-be/service"
	"github.com/quangdnptit/booking-be/storage"
)

func main() {
	ctx := context.Background()

	// ZeroLog
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Info().Str("event", "config_load").Msg(".env not found, using system env")
	}

	// DynamoDB config
	db, err := storage.NewDynamoDBStore(ctx, os.Getenv("DYNAMODB_ENDPOINT"))
	if err != nil {
		log.Fatal().Err(err).Msg("dynamodb client")
	}

	// PostgreSQL
	pgCfg, err := storage.LoadPostgresConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Str("event", "config").Msg("postgres config")
	}
	pgPool, err := storage.NewPostgresPool(ctx, pgCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres client")
	}
	defer pgPool.Close()
	log.Info().Str("event", "postgres_ready").Msg("postgresql pool connected")
	//JWT
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		log.Fatal().Str("event", "config").Msg("JWT_SECRET is required")
	}
	jwtTTL := time.Hour
	if s := os.Getenv("JWT_TTL_SECONDS"); s != "" {
		if sec, err := time.ParseDuration(s + "s"); err == nil && sec > 0 {
			jwtTTL = sec
		}
	}
	refreshTTL := 7 * 24 * time.Hour
	if s := os.Getenv("JWT_REFRESH_TTL_SECONDS"); s != "" {
		if sec, err := time.ParseDuration(s + "s"); err == nil && sec > 0 {
			refreshTTL = sec
		}
	}

	// Init dependencies
	// Repos
	bookingRepo := repo.NewDynamoBookingRepo(db)
	bookedSeatRepo := repo.NewDynamoBookedSeatRepo(db)
	userRepo := repo.NewDynamoUserRepo(db)
	ledgerRepo := repo.NewDynamoLedgerRepo(db)
	movieRepo := repo.NewPostgresProgramRepo(pgPool)

	//services
	bookingSvc := service.NewBookingService(bookingRepo, bookedSeatRepo, userRepo, db)
	seatPriceCfg := service.LoadSeatPriceConfigFromEnv()
	seatService := service.NewSeatService(bookedSeatRepo, seatPriceCfg)
	authSvc := service.NewAuthService(userRepo, jwtSecret, jwtTTL, refreshTTL)
	balanceSvc := service.NewBalanceService(userRepo, ledgerRepo)
	movieSvc := service.NewMovieService(movieRepo)
	balanceHandler := handlers.NewBalanceHandler(balanceSvc)

	//Hanlders
	handler := handlers.NewHandler()
	seatHandler := handlers.NewSeatHandler(seatService)
	bookingHandler := handlers.NewBookingHandler(bookingSvc)
	movieHandler := handlers.NewMovieHandler(movieSvc)
	authHandler := handlers.NewAuthHandler(authSvc)

	// Init Gin Router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(observability.TracingMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Router mapping — public
	router.GET("/api/v1/health", handler.HealthCheck)
	router.POST("/api/v1/auth/login", authHandler.Login)
	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/refresh", authHandler.Refresh)
	router.POST("/api/v1/register", authHandler.Register)
	router.GET("/showtimes/:showtimeId/seats", seatHandler.GetSeats)
	router.GET("/api/movies", movieHandler.ListMovies)
	router.GET("/api/movies/:id", movieHandler.GetMovieByID)
	router.GET("/api/showtimes", movieHandler.ListShowtimes)
	router.GET("/api/showtimes/:id", movieHandler.GetShowtimeByID)
	router.GET("/api/theaters", movieHandler.ListTheaters)
	router.GET("/api/theaters/:id", movieHandler.GetTheaterByID)
	router.GET("/api/rooms/theater/:theaterId", movieHandler.ListRoomsByTheater)

	// JWT Auth middleware config
	protected := router.Group("")
	protected.Use(auth.JWTAuthMiddleware(jwtSecret))
	protected.POST("/api/v1/seats/generate-seats", seatHandler.GenerateSeats)
	protected.GET("/api/v1/users/me", authHandler.GetUserInfo)
	protected.POST("/api/v1/bookings", bookingHandler.BookSeats)
	protected.POST("/api/v1/deposit", balanceHandler.Deposit)
	protected.GET("/api/v1/users/:userId/bookings", bookingHandler.GetUserBookingHistory)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8888"
	} else if port[0] != ':' {
		port = ":" + port
	}
	log.Info().Str("event", "server_listen").Str("addr", port).Msg("listening")
	if err := router.Run(port); err != nil {
		log.Fatal().Err(err).Msg("server")
	}
}
