package api

import (
	"EffectiveMobileAPI/internal/api/request"
	"EffectiveMobileAPI/internal/storage/postgres"
	"encoding/json"
	"github.com/lib/pq"
	"log/slog"
	"net/http"
	"strconv"
)

func AddSongHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.api.AddSongHandler()"

		var song postgres.Song
		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&song)
		if err != nil {
			log.Error("Error decoding request body", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error decoding request body"))
			return
		}

		status, err := storage.AddSong(song, log)
		if err != nil {
			log.Error("Error adding song", "error", err, "status", status)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
			return
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(request.Ok())
	}
}

func ChangeInfoSongHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.api.AddInfoSongHandler()"

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("no id or transmitted incorrectly", "error", err, "id", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var infoSong postgres.InfoSong
		err = json.NewDecoder(r.Body).Decode(&infoSong)
		if err != nil {
			log.Error("Error decoding request body", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error decoding request body"))
			return
		}

		status, err := storage.ChangeInfo(id, infoSong, log)
		if err != nil {
			log.Error("Error changing song info", "error", err, "status", status)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
			return
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(request.Ok())
	}
}
