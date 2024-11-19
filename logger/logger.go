package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLog() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).With().Timestamp().Caller().Logger()
}

func Fatal() *zerolog.Event {
	return Logger.WithLevel(zerolog.FatalLevel)
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}
