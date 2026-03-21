package logger

import (
	"time"

	"github.com/rs/zerolog"
)

type ConsoleWriterOption func(w *zerolog.ConsoleWriter)

func WithTimeFormat(format string) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = format
	}
}

func WithPartsOrder(parts []string) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.PartsOrder = parts
	}
}

func WithNoColor(noColor bool) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.NoColor = noColor
	}
}

func WithTimeLocation(location *time.Location) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.TimeLocation = location
	}
}

func WithPartsExclude(parts []string) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.PartsExclude = parts
	}
}

func WithFieldsOrder(fields []string) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.FieldsOrder = fields
	}
}

func WithFieldsExclude(fields []string) ConsoleWriterOption {
	return func(w *zerolog.ConsoleWriter) {
		w.FieldsExclude = fields
	}
}
