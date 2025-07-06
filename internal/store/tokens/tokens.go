package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Hash      []byte    `json:"hash"`
	UserID    int       `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"scope"`
	PlainText string    `json:"token"`
}

type TokenStore interface {
	Insert(token *Token) error
	Delete(hash []byte) error
	Find(hash []byte) (*Token, error)
}

func GetTokenStore(userID int, scope string, expiry time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Scope:  scope,
		Expiry: time.Now().Add(expiry),
	}

	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
