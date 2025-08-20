package store

import (
	"database/sql"
)

type PostgresCronJobStore struct {
	db *sql.DB
}

func NewPostgresCronJobStore(db *sql.DB) *PostgresCronJobStore {
	return &PostgresCronJobStore{
		db: db,
	}
}

type CronJobStore interface {
	DeleteExpiredTokens() error
}

func (pg *PostgresCronJobStore) DeleteExpiredTokens() error {
	query := `DELETE FROM tokens WHERE expiry < NOW()`
	_, err := pg.db.Exec(query)
	return err
}
