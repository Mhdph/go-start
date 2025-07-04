package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host:localhost user:postgres password:postgres dbname:postgres sslmode:disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Printf("connected to database")

	return db, nil
}
