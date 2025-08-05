package store

import (
	"database/sql"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenreStore interface {
	CreateGenre(*Genre) (*Genre, error)
	GetGenreById(id int) (*Genre, error)
	UpdateGenre(*Genre) error
	DeleteGenreById(id int) error
}
type PostgresGenreStore struct {
	db *sql.DB
}

func NewPostgresGenreStore(db *sql.DB) *PostgresGenreStore {
	return &PostgresGenreStore{db: db}
}

func (pg *PostgresGenreStore) CreateGenre(genre *Genre) (*Genre, error) {
	query := `
	INSERT INTO genres (id, name)
	VALUES ($1, $2)
	RETURNING id
	`

	err := pg.db.QueryRow(query, genre.ID, genre.Name).Scan(&genre.ID)

	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (pg *PostgresGenreStore) GetGenreById(id int) (*Genre, error) {
	query := `
	SELECT id, name
	FROM genres
	WHERE id = $1
	`

	genre := &Genre{}
	err := pg.db.QueryRow(query, id).Scan(&genre.ID, &genre.Name)

	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (pg *PostgresGenreStore) DeleteGenreById(id int) error {
	query := `
	DELETE FROM genres
	WHERE id = $1
	`

	_, err := pg.db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
