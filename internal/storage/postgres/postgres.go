package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"time"
)

type Song struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type InfoSong struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddSong(song Song, log *slog.Logger) (int, error) {

	const op = "storage.postgres.AddSong()"

	query := `INSERT INTO song (name, music_group) VALUES ($1, $2)`

	_, err := s.db.Exec(query, song.Name, song.Group)
	if err != nil {
		log.Error("Error to insert", op)
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
