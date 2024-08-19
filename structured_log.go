package main

import (
	"log/slog"
	"os"
	"time"
)

func structuredLog() {
	slog.Error("error log msg")
	slog.Warn("warn log msg")
	slog.Info("info log msg 2")
	slog.Debug("debug log msg")

	userID := "fred"
	loginCount := 20
	slog.Info("user login", "id", userID, "login_count", loginCount)

	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(os.Stderr, options)
	mySlog := slog.New(handler)
	lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
	mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)
}
