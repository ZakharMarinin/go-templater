package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func Logger(env string, path string) (*slog.Logger, error) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	logFile, err := os.OpenFile(
		filepath.Join(path, "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	
	var log *slog.Logger

	switch env {
	case EnvLocal:
		log = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, fmt.Errorf("unknown env")
	}

	return log, nil
}