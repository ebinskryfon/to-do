package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New builds a zerolog.Logger: pretty console output in development,
// structured JSON in any other environment (e.g. production).
func New(env string) zerolog.Logger {
	if env == "development" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
