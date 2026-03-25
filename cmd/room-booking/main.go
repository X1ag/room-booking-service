package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/booking"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/http/handlers"
	"test-backend-1-X1ag/internal/http/middleware"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/repository/postgres"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/schedule"
	"test-backend-1-X1ag/internal/slot"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := logger.NewZerologLogger(cfg.Logger)
	if err != nil {
		log.Panicln("cant create logger, err: ", err)
		return
	}

	// Init all loggers 
	appLogger := logger.WithFeature("app")
	slotLogger := logger.WithFeature("slot")
	roomLogger := logger.WithFeature("room")
	scheduleLogger := logger.WithFeature("schedule")
	bookingLogger := logger.WithFeature("booking")

	appLogger.Info().Msg("Loggers with features created")

	appLogger.Info().Msg("Connecting to database...")
	pool, err := postgres.ConnectPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	appLogger.Info().Msg("Connected to database")

	// migrations
	m, err := migrate.New(
		"file://migrations",
		cfg.DB.DSN(),
	)
	if err != nil {
		log.Fatalf("create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("run migrations: %v", err)
	}

	appLogger.Info().Msg("Migrations applied successfully")
	// TODO: add migrations from file://migrations
	
	// init repositories
	slotRepo := postgres.NewSlotRepository(pool)
	roomRepo := postgres.NewRoomRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	bookingRepo := postgres.NewBookingRepository(pool)

	// init usecases
	slotUsecase := slot.NewSlotUsecase(slotRepo, roomRepo, scheduleRepo, slotLogger)
	roomUsecase := room.NewRoomUsecase(roomRepo, roomLogger)
	scheduleUsecase := schedule.NewSheduleUsecase(scheduleRepo, roomRepo, scheduleLogger)
	bookingUsecase := booking.NewBookingUsecase(bookingRepo, bookingLogger)

	jwtManager := auth.NewJWTManager(cfg.Auth)
	authUsecase := auth.NewAuthUsecase(jwtManager, cfg.Auth, logger)

	// init handlers
	slotHandlers := handlers.NewSlotHandler(slotUsecase)
	roomHandlers := handlers.NewRoomHandler(roomUsecase)
	scheduleHandlers := handlers.NewScheduleHandler(scheduleUsecase)
	bookingHandlers := handlers.NewBookingHandler(bookingUsecase)
	authHandler := handlers.NewAuthHandler(authUsecase)
	_ = slotHandlers
	_ = bookingHandlers


	// init server
	r := gin.Default()

	// init middleware
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(jwtManager, appLogger))

	admin := authorized.Group("/")
	admin.Use(middleware.RequireRole("admin"))

	user := authorized.Group("/")
	user.Use(middleware.RequireRole("user"))

	// init routes
	admin.POST("/rooms/create", roomHandlers.Create())	
	admin.POST("/rooms/:roomId/schedule/create", scheduleHandlers.Create())
	admin.GET("/bookings/list", bookingHandlers.GetUserBookings())
	
	authorized.GET("/rooms/list", roomHandlers.GetRooms())
	authorized.GET("/rooms/:roomId/slots/list", slotHandlers.GetSlotsByRoomID())

	user.POST("/bookings/create", bookingHandlers.Create())
	user.POST("/bookings/:bookingId/cancel", bookingHandlers.Cancel())
	
	r.Handle("GET", "/_info", handlers.Info)
	r.Handle("POST", "/dummyLogin", authHandler.DummyLogin)

	if err := r.Run(cfg.HTTP.Addr()); err != nil {
		log.Fatalf("run HTTP server: %v", err)
	}
}
