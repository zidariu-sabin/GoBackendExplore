package store

import (
	"GoBackendExploreMovieTracker/internal/tokens"
	"database/sql"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgrestokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	CreateNewToken(userID int64, expiry time.Duration, scope string) (*tokens.Token, error)
	InsertToken(token *tokens.Token) error
}

func (pg *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = pg.InsertToken(token)

	return token, err
}

func (pg *PostgresTokenStore) InsertToken(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	_, err := pg.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)

	return err
}
