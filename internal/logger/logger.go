package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"test-backend-1-X1ag/internal/config"

	"github.com/rs/zerolog"
)

type ZerologLogger struct {
	logger zerolog.Logger
	cfg    config.LoggerConfig
}

func NewZerologLogger(cfg config.LoggerConfig, opts ...ConsoleWriterOption) (*ZerologLogger, error) {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	// Создаём директорию логов
	if err := os.MkdirAll(cfg.LogsDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Создаём имя лог-файла: app_2025-04-09.log
	date := time.Now().Format(time.DateOnly)
	filename := filepath.Join(cfg.LogsDir, "app_"+date+".log")

	// Открываем файл
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Создаём writer на консоль
	consoleWriter := zerolog.NewConsoleWriter()
	for _, opt := range opts {
		opt(&consoleWriter)
	}

	// Комбинируем вывод в файл + консоль
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	// Настраиваем логгер
	logger := zerolog.New(multi).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &ZerologLogger{
		logger: logger,
		cfg:    cfg,
	}, nil
}

func NewTestLogger() *ZerologLogger {
	logger := zerolog.New(io.Discard).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()

	return &ZerologLogger{
		logger: logger,
		cfg:    config.LoggerConfig{},
	}
}

// Создание "именованного" логгера — для фичи
func (l *ZerologLogger) WithFeature(feature string) *ZerologLogger {
	child := l.logger.With().Str("feature", feature).Logger()
	return &ZerologLogger{
		logger: child,
		cfg:    l.cfg,
	}
}

// Ниже просто проксируем уровни
func (l *ZerologLogger) Trace() *zerolog.Event { return l.logger.Trace() }
func (l *ZerologLogger) Debug() *zerolog.Event { return l.logger.Debug() }
func (l *ZerologLogger) Info() *zerolog.Event  { return l.logger.Info() }
func (l *ZerologLogger) Warn() *zerolog.Event  { return l.logger.Warn() }
func (l *ZerologLogger) Error() *zerolog.Event { return l.logger.Error() }
func (l *ZerologLogger) Fatal() *zerolog.Event { return l.logger.Fatal() }
func (l *ZerologLogger) Panic() *zerolog.Event { return l.logger.Panic() }
func (l *ZerologLogger) Log() *zerolog.Event   { return l.logger.Log() }
func (l *ZerologLogger) With() zerolog.Context { return l.logger.With() }
func (l *ZerologLogger) Level(lvl zerolog.Level) zerolog.Logger {
	return l.logger.Level(lvl)
}
func (l *ZerologLogger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}
func (l *ZerologLogger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}
