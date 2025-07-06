package store

import (
	"database/sql"
	"time"

	"github.com/mhdph/go-start/internal/store/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	Create(userID string, scope string, expiry time.Duration) (*tokens.Token, error)
}
