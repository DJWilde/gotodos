package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
	"todos/internal/validator"
)

// Token struct and functions
type Token struct {
	Plaintext 	string 		`json:"token"`
	Hash 		[]byte 		`json:"-"`
	UserID		int64		`json:"-"`
	Expiry 		time.Time	`json:"expiry"`
}

func GenerateToken(userID int64, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:] // conversion of hash to slice

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// DB model

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration) (*Token, error) {
	token, err := GenerateToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry)
		VALUES ($1, $2, $3)`

	args := []any{token.Hash, token.UserID, token.Expiry}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}