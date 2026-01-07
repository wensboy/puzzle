package clog

import (
	"log/slog"
	"os"
	"testing"
)

func TestPlainHandler(t *testing.T) {
	logger := slog.New(NewPlainHandler(&LogOption{
		Out:   os.Stderr,
		Level: slog.LevelInfo,
	}))
	t.Run("plain handler - debug", func(t *testing.T) {
		logger.Debug("plain handler debug test...pass")
	})
	t.Run("plain handler - info", func(t *testing.T) {
		logger.Info("plain handler info test...pass")
	})
	t.Run("plain handler - warn", func(t *testing.T) {
		logger.Warn("plain handler warn test...pass")
	})
}
