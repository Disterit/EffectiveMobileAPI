package storage

import (
	"EffectiveMobileAPI/internal/config"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func Connection(log *slog.Logger) *sql.DB {

	const op = "storage.connection.Connection()"

	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Error("Error to connect database", err, op)
	}

	return db
}
