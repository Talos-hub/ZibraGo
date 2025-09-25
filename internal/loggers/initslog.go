package loggers

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const LogPath = "logs/log.log"

func SetupLogger() (*slog.Logger, error) {
	err := createLogFolder()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})), nil
}

// Create log directory if it doesn't exist
func createLogFolder() error {
	dir := filepath.Dir(LogPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	return nil
}
