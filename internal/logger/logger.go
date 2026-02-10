package logger

import (
	"log/slog"
	"os"
)

func Init() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("logger initialized")
}
