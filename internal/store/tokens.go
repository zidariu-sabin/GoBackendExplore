package store

import (
	"GoBackendExploreMovieTracker/internal/tokens"
	"database/sql"
	"fmt"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	CreateNewToken(userID int64, expiry time.Duration, scope string) (*tokens.Token, error)
	InsertToken(token *tokens.Token) error
	DeleteOldTokens() (int64, error)
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

func (pg *PostgresTokenStore) DeleteOldTokens() (int64, error) {
	query := `
	DELETE FROM tokens
	WHERE expiry > NOW ()
	RETURNING count
	`

	result, err := pg.db.Exec(query)

	fmt.Printf("Deleted rows: %v", result)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
