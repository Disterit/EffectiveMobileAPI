package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"time"
)

type Library struct {
	Songs Songs `json:"songs"`
}

type Songs struct {
	Song     Song     `json:"song"`
	InfoSong InfoSong `json:"info_song"`
}

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

	return id, nil
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

func (s *Storage) DeleteSong(id int, log *slog.Logger) (sql.Result, error) {
	const op = "storage.postgres.DeleteInfo()"

	query := `DELETE FROM Song WHERE id = $1;`

	res, err := s.db.Exec(query, id)
	if err != nil {
		log.Error("Error to delete", op)
	}

	return res, nil
}

func (s *Storage) GetText(id int, log *slog.Logger) (string, error) {
	const op = "storage.postgres.GetText()"

	query := `SELECT text FROM infosong WHERE id_song = $1;`

	var text string

	err := s.db.QueryRow(query, id).Scan(&text)
	if err != nil {
		log.Error("Error to get song text", op)
	}

	return text, err
}

func (s *Storage) GetLibrary(log *slog.Logger) ([]Library, error) {

	const op = "storage.postgres.GetLibrary()"

	query := `SELECT s.id, s.music_group, s.name, i.text, i.releasedate, i.link
				FROM song s
				JOIN infosong i ON s.id = i.id_song;
				`

	var library []Library

	rows, err := s.db.Query(query)
	if err != nil {
		log.Error("Error to get songs", op)
	}

	for rows.Next() {
		var lib Library
		var id int64
		err = rows.Scan(&id,
			&lib.Songs.Song.Group,
			&lib.Songs.Song.Name,
			&lib.Songs.InfoSong.Text,
			&lib.Songs.InfoSong.ReleaseDate,
			&lib.Songs.InfoSong.Link)
		if err != nil {
			log.Error("Error to get songs", op)
			return nil, err
		}

		library = append(library, lib)
	}

	return library, nil
}

func (s *Storage) GetInfo(name, group string, log *slog.Logger) (InfoSong, error) {

	const op = "storage.postgres.GetInfo()"

	query := `SELECT text, releasedate, link FROM Library WHERE music_group = $1 AND name = $2;`

	var infoSong InfoSong

	rows, err := s.db.Query(query, name, group)

	if err != nil {
		log.Error("Error to get songs", op)
	}

	for rows.Next() {
		err = rows.Scan(&infoSong.Text,
			&infoSong.ReleaseDate,
			&infoSong.Link)
		if err != nil {
			log.Error("Error to get songs", op)
			return InfoSong{}, err
		}
	}

	return infoSong, nil
}
