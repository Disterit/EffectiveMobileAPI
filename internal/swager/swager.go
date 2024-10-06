package swager

import (
	"EffectiveMobileAPI/internal/api"
	"EffectiveMobileAPI/internal/storage/postgres"
	"github.com/go-chi/chi"
	"log/slog"
)

func InitRoutes(r *chi.Mux, log *slog.Logger, storage *postgres.Storage) {
	// Ваши хендлеры для API
	r.HandleFunc("/songs/add", api.AddSongHandler(log, storage))
	r.HandleFunc("/songs/change", api.ChangeInfoSongHandler(log, storage))
	r.HandleFunc("/songs/delete", api.DeleteSongHandler(log, storage))
	r.HandleFunc("/songs/text", api.TextSongHandler(log, storage))
	r.HandleFunc("/library", api.LibraryHandler(log, storage))
	r.HandleFunc("/info", api.InfoHandler(log, storage))
}
