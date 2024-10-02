package main

import (
	"EffectiveMobileAPI/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	encDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting api", slog.String("key", cfg.Env))
	log.Debug("debug message enable")

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case encDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
