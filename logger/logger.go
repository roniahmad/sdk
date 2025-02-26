package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

func NewLog(logTpye string) (logger zerolog.Logger) {
	zerolog.TimeFieldFormat = time.RFC3339
	logger = zerolog.New(os.Stdout).With().Timestamp().Str("role", filepath.Base(os.Args[0])).Logger()

	if logTpye == "console" {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
		logger = zerolog.New(output).With().Timestamp().Str("role", filepath.Base(os.Args[0])).Logger()
	}

	return
}
