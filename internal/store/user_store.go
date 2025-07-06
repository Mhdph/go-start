package store

import (
	"database/sql"
	"fmt"
	"time"
)

type password struct {
	plainText *string
	hash      []byte
}
type User struct {
	ID        int       `json:"Id"`
	Username  string    `json:"string"`
	Email     string    `json:"email"`
	Password  password  `json:"_"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password, bio) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	row := s.db.QueryRow(query, user.Username, user.Email, user.Password, user.Bio)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := ` 
	UPDATE users 
	SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP 
	WHERE id= $4 
	RETURNING updated_at
	`

	result, err := s.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		Password: password{},
	}

	query := ` 
	SELECT id, username, email, bio, created_at, updated_at 
	FROM users
	WHERE username = $1 
	`
	err := s.db.QueryRow(query, username).Scan(
		&user.Username,
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}
