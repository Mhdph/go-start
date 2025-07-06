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
	Create(userID int, scope string, expiry time.Duration) (*tokens.Token, error)
}

func (t *PostgresTokenStore) CreateNewToken(userID int, scope string, expiry time.Duration) (*tokens.Token, error) {
	token, err := tokens.GetTokenStore(userID, scope, expiry)

	if err != nil {
		return nil, err
	}

	err = t.Insert(token)

	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (id, user_id, scope, expiry)
		VALUES ($1, $2, $3, $4)
	`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Scope, token.Expiry)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(UserID string, scope string) error {
	query := ` 
	DELETE FROM tokens 
	WHERE scope = $1 AND user_id = $2 
	`
	_, err := t.db.Exec(query, scope, UserID)
	return err
}
