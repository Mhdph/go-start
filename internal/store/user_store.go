package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.plainText = &plainText
	p.hash = hash
	return nil
}

func (p password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
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

var AnonymousUser = &User{
	ID:       0,
	Username: "anonymous",
	Email:    "",
	Bio:      "",
	// Password and timestamps can be left zero value
}

func (u *User) IsAnnoymous() bool {
	return u == AnonymousUser
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
	GetUserToken(scope, tokenPlainText string) (*User, error)
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

func (s *PostgresTokenStore) GetUserToken(scope, plainTextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	query := `
	  SELECT u.id, u.email, u.username, u.password, u.bio, u.created_at, u.updated_at 
	  FROM users u 	
	  INNER JOIN tokens t ON t.users.id = u.user.id 
	  WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3 
	`

	user := &User{
		Password: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err == nil {
		return nil, err
	}

	return user, nil
}
func (s *PostgresUserStore) GetUserToken(scope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
	SELECT u.id, u.email, u.username, u.password, u.bio, u.created_at, u.updated_at 
	FROM users u 
	INNER JOIN tokens t ON t.users.id = u.user.id 
	WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3`

	user := &User{
		Password: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get user token: %w", err)
	}

	return user, nil
}
