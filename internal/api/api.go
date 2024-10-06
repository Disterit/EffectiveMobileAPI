package api

import (
	"EffectiveMobileAPI/internal/api/request"
	"EffectiveMobileAPI/internal/api/response"
	"EffectiveMobileAPI/internal/storage/postgres"
	"encoding/json"
	"fmt"
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

		url := fmt.Sprintf("http://0.0.0.0:8081/info?group=%s&song=%s", song.Group, song.Name)
		infoSong, err := response.GetInfoSong(log, url)
		if err != nil {
			log.Error("Error getting info song in library", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error decoding request body"))
			return
		}

		id, err := storage.AddSong(song, log)
		if err != nil {
			log.Error("Error adding song", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
			return
		}

		err = response.ChangeInfoSong(log, infoSong, id)
		if err != nil {
			log.Error("Error adding info song", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(request.Ok())
	}
}

func ChangeInfoSongHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.AddInfoSongHandler()"

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("no id or transmitted incorrectly", "error", err, "operation", op)
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
			log.Error("Error changing song info", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
			return
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(request.Ok())
	}
}

func DeleteSongHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.DeleteSongHandler()"

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("no id or transmitted incorrectly", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := storage.DeleteSong(id, log)
		if err != nil {
			log.Error("Error deleting song", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
		}

		rowsAffected, err := result.RowsAffected()
		if rowsAffected == 0 || err != nil {
			log.Error("Error deleting song, song with this id not found", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error deleting song, song id not found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(request.Ok())
	}
}

func TextSongHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.TextSongHandler()"

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("no id or transmitted incorrectly", "error", err, "operation", op)
			return
		}

		text, err := storage.GetText(id, log)
		if err != nil {
			log.Error("Error getting song text", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			pgErr, _ := err.(*pq.Error)
			json.NewEncoder(w).Encode(request.BadRequest(pgErr.Message))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(text)
	}
}

func LibraryHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.LibraryHandler()"

		library, err := storage.GetLibrary(log)
		if err != nil {
			log.Error("Error getting library", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error getting library"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library)
		w.WriteHeader(http.StatusOK)
	}
}

func InfoHandler(log *slog.Logger, storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.InfoHandler()"

		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")

		info, err := storage.GetInfo(group, song, log)
		if err != nil {
			log.Error("Error getting info", "error", err, "operation", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(request.BadRequest("Error getting info"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)

	}
}
