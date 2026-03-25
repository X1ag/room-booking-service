package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	App        AppConfig
	HTTP       HTTPConfig
	DB         DBConfig
	Auth       AuthConfig
	Migrations MigrationConfig
	Logger     LoggerConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int32
	MinIdleConns    int32
	MaxConnLifetime time.Duration
}

type AuthConfig struct {
	JWTSecret    string
	TokenTTL     time.Duration
	DummyAdminID uuid.UUID
	DummyUserID  uuid.UUID
}

type LoggerConfig struct {
	LogLevel string `env:"LOGGER_LEVEL" envDefault:"debug"`
	NoColor  bool   `env:"LOGGER_NO_COLOR" envDefault:"false"`

	TimeFormat   string `env:"LOGGER_TIME_FORMAT" envDefault:"2006-01-02T15:04:05Z"`
	TimeLocation string `env:"LOGGER_TIME_LOCATION" envDefault:"UTC"`

	PartsOrder    string `env:"LOGGER_PARTS_ORDER" envDefault:"time,level,logger,message"`
	PartsExclude  string `env:"LOGGER_PARTS_EXCLUDE" envDefault:""`
	FieldsOrder   string `env:"LOGGER_FIELDS_ORDER" envDefault:""`
	FieldsExclude string `env:"LOGGER_FIELDS_EXCLUDE" envDefault:""`

	LogsDir string `env:"LOGGER_LOGS_DIR" envDefault:"./logs"`
}

type MigrationConfig struct {
	Path string
}

func Load() (Config, error) {
	httpPort, err := getInt("HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}

	httpReadTimeout, err := getDuration("HTTP_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}

	httpWriteTimeout, err := getDuration("HTTP_WRITE_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}

	httpShutdownTimeout, err := getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}

	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return Config{}, err
	}

	dbMaxOpenConns, err := getInt("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	dbMinIdleConns, err := getInt("DB_MIN_IDLE_CONNS", 2)
	if err != nil {
		return Config{}, err
	}

	dbMaxConnLifetime, err := getDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}

	tokenTTL, err := getDuration("JWT_TOKEN_TTL", 24*time.Hour)
	if err != nil {
		return Config{}, err
	}

	loggerNoColor, err := getBool("LOGGER_NO_COLOR", false)
	if err != nil {
		return Config{}, err
	}

	adminUUID, err := getUUID("DUMMY_ADMIN_ID", uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	if err != nil {
		return Config{}, err
	}
	userUUID, err := getUUID("DUMMY_USER_ID", uuid.MustParse("22222222-2222-2222-2222-222222222222"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App: AppConfig{
			Name: getString("APP_NAME", "room-booking"),
			Env:  getString("APP_ENV", "local"),
		},
		HTTP: HTTPConfig{
			Host:            getString("HTTP_HOST", "0.0.0.0"),
			Port:            httpPort,
			ReadTimeout:     httpReadTimeout,
			WriteTimeout:    httpWriteTimeout,
			ShutdownTimeout: httpShutdownTimeout,
		},
		DB: DBConfig{
			Host:            getString("DB_HOST", "postgres"),
			Port:            dbPort,
			Name:            getString("DB_NAME", "room_booking"),
			User:            getString("DB_USER", "postgres"),
			Password:        getString("DB_PASSWORD", "postgres"),
			SSLMode:         getString("DB_SSLMODE", "disable"),
			MaxOpenConns:    int32(dbMaxOpenConns),
			MinIdleConns:    int32(dbMinIdleConns),
			MaxConnLifetime: dbMaxConnLifetime,
		},
		Auth: AuthConfig{
			JWTSecret:    getString("JWT_SECRET", "local-dev-secret-change-me"),
			TokenTTL:     tokenTTL,
			DummyAdminID: adminUUID,
			DummyUserID:  userUUID,
		},
		Migrations: MigrationConfig{
			Path: getString("MIGRATIONS_PATH", "file://migrations"),
		},
		Logger: LoggerConfig{
			LogLevel:      getString("LOGGER_LEVEL", "debug"),
			NoColor:       loggerNoColor,
			TimeFormat:    getString("LOGGER_TIME_FORMAT", "2006-01-02T15:04:05Z"),
			TimeLocation:  getString("LOGGER_TIME_LOCATION", "UTC"),
			PartsOrder:    getString("LOGGER_PARTS_ORDER", "time,level,logger,message"),
			PartsExclude:  getString("LOGGER_PARTS_EXCLUDE", ""),
			FieldsOrder:   getString("LOGGER_FIELDS_ORDER", ""),
			FieldsExclude: getString("LOGGER_FIELDS_EXCLUDE", ""),
			LogsDir:       getString("LOGGER_LOGS_DIR", "./logs"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var errs []error

	if strings.TrimSpace(c.App.Name) == "" {
		errs = append(errs, errors.New("APP_NAME must not be empty"))
	}
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		errs = append(errs, fmt.Errorf("invalid HTTP_PORT: %d", c.HTTP.Port))
	}
	if c.HTTP.ReadTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_READ_TIMEOUT must be positive"))
	}
	if c.HTTP.WriteTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_WRITE_TIMEOUT must be positive"))
	}
	if c.HTTP.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_SHUTDOWN_TIMEOUT must be positive"))
	}
	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		errs = append(errs, errors.New("JWT_SECRET must not be empty"))
	}
	if c.Auth.TokenTTL <= 0 {
		errs = append(errs, errors.New("JWT_TOKEN_TTL must be positive"))
	}
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		errs = append(errs, fmt.Errorf("invalid DB_PORT: %d", c.DB.Port))
	}
	if strings.TrimSpace(c.DB.Host) == "" {
		errs = append(errs, errors.New("DB_HOST must not be empty"))
	}
	if strings.TrimSpace(c.DB.Name) == "" {
		errs = append(errs, errors.New("DB_NAME must not be empty"))
	}
	if strings.TrimSpace(c.DB.User) == "" {
		errs = append(errs, errors.New("DB_USER must not be empty"))
	}
	if strings.TrimSpace(c.DB.SSLMode) == "" {
		errs = append(errs, errors.New("DB_SSLMODE must not be empty"))
	}
	if c.DB.MaxOpenConns <= 0 {
		errs = append(errs, errors.New("DB_MAX_OPEN_CONNS must be positive"))
	}
	if c.DB.MinIdleConns < 0 {
		errs = append(errs, errors.New("DB_MIN_IDLE_CONNS must not be negative"))
	}
	if c.DB.MaxConnLifetime <= 0 {
		errs = append(errs, errors.New("DB_MAX_CONN_LIFETIME must be positive"))
	}
	if c.Auth.DummyAdminID == uuid.Nil {
		errs = append(errs, errors.New("DUMMY_ADMIN_USER_ID must not be empty"))
	}
	if c.Auth.DummyUserID == uuid.Nil {
		errs = append(errs, errors.New("DUMMY_USER_USER_ID must not be empty"))
	}
	if strings.TrimSpace(c.Migrations.Path) == "" {
		errs = append(errs, errors.New("MIGRATIONS_PATH must not be empty"))
	}
	if strings.TrimSpace(c.Logger.LogLevel) == "" {
		errs = append(errs, errors.New("LOGGER_LEVEL must not be empty"))
	}
	if !isValidLogLevel(c.Logger.LogLevel) {
		errs = append(errs, fmt.Errorf("invalid LOGGER_LEVEL: %s", c.Logger.LogLevel))
	}
	if strings.TrimSpace(c.Logger.TimeFormat) == "" {
		errs = append(errs, errors.New("LOGGER_TIME_FORMAT must not be empty"))
	}
	if strings.TrimSpace(c.Logger.TimeLocation) == "" {
		errs = append(errs, errors.New("LOGGER_TIME_LOCATION must not be empty"))
	} else if _, err := time.LoadLocation(c.Logger.TimeLocation); err != nil {
		errs = append(errs, fmt.Errorf("invalid LOGGER_TIME_LOCATION: %w", err))
	}
	if strings.TrimSpace(c.Logger.LogsDir) == "" {
		errs = append(errs, errors.New("LOGGER_LOGS_DIR must not be empty"))
	}

	return errors.Join(errs...)
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be int: %w", key, err)
	}

	return parsed, nil
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be duration: %w", key, err)
	}

	return parsed, nil
}

func getBool(key string, fallback bool) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be bool: %w", key, err)
	}

	return parsed, nil
}

func getUUID(key string, fallback uuid.UUID) (uuid.UUID, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s must be uuid: %w", key, err)
	}

	return parsed, nil
}

func isValidLogLevel(level string) bool {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace", "debug", "info", "warn", "error", "fatal", "panic":
		return true
	default:
		return false
	}
}
