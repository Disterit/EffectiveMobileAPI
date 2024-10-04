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
	ReleaseDate *time.Time `json:"releaseDate"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddSong(song Song, log *slog.Logger) (int, error) {

	const op = "storage.postgres.AddSong()"

	query := `INSERT INTO song (name, music_group) VALUES ($1, $2) returning id`

	var id int

	err := s.db.QueryRow(query, song.Name, song.Group).Scan(&id)
	if err != nil {
		log.Error("Error to insert", op)
		return http.StatusBadRequest, err
	}

	query = `INSERT INTO infosong (id_song) VALUES ($1)`

	_, err = s.db.Exec(query, id)
	if err != nil {
		log.Error("Error to insert", op)
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (s *Storage) ChangeInfo(id int, info InfoSong, log *slog.Logger) (int, error) {

	const op = "storage.postgres.AddInfo()"

	query := `
		UPDATE InfoSong
		SET 
		    releaseDate = COALESCE($1, releaseDate),
		    text = COALESCE($2, text),
		    link = COALESCE($3, link)
		WHERE id_song = $4;
	`

	_, err := s.db.Exec(query, info.ReleaseDate, info.Text, info.Link, id)
	if err != nil {
		log.Error("Error to update", op)
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
