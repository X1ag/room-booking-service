package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/booking"
	"test-backend-1-X1ag/internal/conference"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/http/handlers"
	"test-backend-1-X1ag/internal/http/middleware"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/repository/postgres"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/schedule"
	"test-backend-1-X1ag/internal/slot"
)

type App struct {
	Router *gin.Engine
	Pool   *pgxpool.Pool
	Logger *logger.ZerologLogger
}

func New(ctx context.Context, cfg config.Config, baseLogger *logger.ZerologLogger) (*App, error) {
	appLogger := baseLogger.WithFeature("app")
	slotLogger := baseLogger.WithFeature("slot")
	roomLogger := baseLogger.WithFeature("room")
	scheduleLogger := baseLogger.WithFeature("schedule")
	bookingLogger := baseLogger.WithFeature("booking")

	appLogger.Info().Msg("Loggers with features created")
	appLogger.Info().Msg("Connecting to database...")

	pool, err := postgres.ConnectPool(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	appLogger.Info().Msg("Connected to database")

	if err := runMigrations(cfg); err != nil {
		pool.Close()
		return nil, err
	}

	appLogger.Info().Msg("Migrations applied successfully")

	slotRepo := postgres.NewSlotRepository(pool)
	roomRepo := postgres.NewRoomRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	bookingRepo := postgres.NewBookingRepository(pool)
	conferenceService := conference.NewMockService()

	slotUsecase := slot.NewSlotUsecase(slotRepo, roomRepo, scheduleRepo, slotLogger)
	roomUsecase := room.NewRoomUsecase(roomRepo, roomLogger)
	scheduleUsecase := schedule.NewSheduleUsecase(scheduleRepo, roomRepo, scheduleLogger)
	bookingUsecase := booking.NewBookingUsecase(bookingRepo, slotRepo, conferenceService, bookingLogger)

	jwtManager := auth.NewJWTManager(cfg.Auth)
	authUsecase := auth.NewAuthUsecase(jwtManager, cfg.Auth, baseLogger)

	slotHandlers := handlers.NewSlotHandler(slotUsecase)
	roomHandlers := handlers.NewRoomHandler(roomUsecase)
	scheduleHandlers := handlers.NewScheduleHandler(scheduleUsecase)
	bookingHandlers := handlers.NewBookingHandler(bookingUsecase)
	authHandler := handlers.NewAuthHandler(authUsecase)

	router := NewRouter(jwtManager, appLogger, roomHandlers, scheduleHandlers, slotHandlers, bookingHandlers, authHandler)

	return &App{
		Router: router,
		Pool:   pool,
		Logger: appLogger,
	}, nil
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
}

func NewRouter(
	jwtManager auth.TokenManager,
	appLogger *logger.ZerologLogger,
	roomHandlers *handlers.RoomHandler,
	scheduleHandlers *handlers.ScheduleHandler,
	slotHandlers *handlers.SlotHandler,
	bookingHandlers *handlers.BookingHandler,
	authHandler *handlers.AuthHandler,
) *gin.Engine {
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(jwtManager, appLogger))

	admin := authorized.Group("/")
	admin.Use(middleware.RequireRole("admin"))

	user := authorized.Group("/")
	user.Use(middleware.RequireRole("user"))

	admin.POST("/rooms/create", roomHandlers.Create())
	admin.POST("/rooms/:roomId/schedule/create", scheduleHandlers.Create())
	admin.GET("/bookings/list", bookingHandlers.ListBookings())

	authorized.GET("/rooms/list", roomHandlers.GetRooms())
	authorized.GET("/rooms/:roomId/slots/list", slotHandlers.GetSlotsByRoomID())

	user.POST("/bookings/create", bookingHandlers.Create())
	user.POST("/bookings/:bookingId/cancel", bookingHandlers.Cancel())
	user.GET("/bookings/my", bookingHandlers.GetUserBookings())

	r.Handle("GET", "/_info", handlers.Info)
	r.Handle("POST", "/dummyLogin", authHandler.DummyLogin)

	return r
}

func runMigrations(cfg config.Config) error {
	m, err := migrate.New(cfg.Migrations.Path, cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			fmt.Printf("failed to close migration source: %v\n", sourceErr)
		}
		if dbErr != nil {
			fmt.Printf("failed to close migration database: %v\n", dbErr)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
