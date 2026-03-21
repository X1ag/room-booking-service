package main

import (
	"context"
	"log"
	"runtime/debug"
	"test-backend-1-X1ag/internal/booking"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/http/handlers"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/repository/postgres"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/schedule"
	"test-backend-1-X1ag/internal/slot"

	"github.com/gin-gonic/gin"
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
	// TODO: add migrations from file://migrations
	
	// init repositories
	slotRepo := postgres.NewSlotRepository(pool)
	roomRepo := postgres.NewRoomRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	bookingRepo := postgres.NewBookingRepository(pool)

	// init usecases
	slotUsecase := slot.NewSlotUsecase(slotRepo, slotLogger)
	roomUsecase := room.NewRoomUsecase(roomRepo, roomLogger)
	scheduleUsecase := schedule.NewSheduleUsecase(scheduleRepo, scheduleLogger)
	bookingUsecase := booking.NewBookingUsecase(bookingRepo, bookingLogger)

	// init middleware
	// TODO: add middleware 

	// init handlers
	slotHandlers := handlers.NewSlotHandler(slotUsecase)
	roomHandlers := handlers.NewRoomHandler(roomUsecase)
	scheduleHandlers := handlers.NewScheduleHandler(scheduleUsecase)
	bookingHandlers := handlers.NewBookingHandler(bookingUsecase)
	_ = slotHandlers
	_ = roomHandlers
	_ = scheduleHandlers
	_ = bookingHandlers

	// init routes

	// init server
	r := gin.Default()

	r.Handle("GET", "/_info", handlers.Info)

	if err := r.Run(cfg.HTTP.Addr()); err != nil {
		log.Fatalf("run HTTP server: %v", err)
	}
}
