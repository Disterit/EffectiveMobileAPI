package main

import (
	"EffectiveMobileAPI/internal/api"
	"EffectiveMobileAPI/internal/config"
	"EffectiveMobileAPI/internal/storage"
	"EffectiveMobileAPI/internal/storage/postgres"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
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

	db := storage.Connection(log)
	router := chi.NewRouter()

	storageDB := postgres.NewStorage(db)

	router.Post("/EffectiveMobile/AddSong", api.AddSongHandler(log, storageDB))

	err := http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Error("Error starting server", err)
	}

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
