package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection.
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

		slog.Info("Successfully connected to the database!")
	return db, nil
}
