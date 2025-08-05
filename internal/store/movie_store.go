package store

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Overview    string    `json:"overview"`
	PosterPath  string    `json:"poster_path"`
	GenreIds    []int     `json:"genre_ids"`
}

type PostgresMovieStore struct {
	db *sql.DB
}

func NewPostgresMovieStore(db *sql.DB) *PostgresMovieStore {
	return &PostgresMovieStore{db: db}
}

type MovieStore interface {
	CreateMovie(*Movie) error
	GetMovieById(id int64) (*Movie, error)
	UpdateMovie(*Movie) error
	DeleteMovieById(id int64) error
}

func (pg *PostgresMovieStore) CreateMovie(movie *Movie) error {

	query := `
	INSERT INTO movies (Title, ReleaseDate, Overview, PosterPath, GenreIds)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	err := pg.db.QueryRow(query, movie.Title, movie.ReleaseDate, movie.Overview, movie.PosterPath, movie.GenreIds).Scan(&movie.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresMovieStore) GetMovieById(id int64) (*Movie, error) {
	query := `
	SELECT id, title, release_date, overview, poster_path, genre_ids
	FROM movies
	WHERE id = $1
	`

	movie := &Movie{}
	err := pg.db.QueryRow(query, id).Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Overview, &movie.PosterPath, &movie.GenreIds)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (pg *PostgresMovieStore) UpdateMovie(movie *Movie) error {
	query := `
	UPDATE movies
	SET title = $1, release_date = $2, overview = $3, poster_path = $4, genre_ids = $5
	WHERE id = $6
	`

	result, err := pg.db.Exec(query, movie.Title, movie.ReleaseDate, movie.Overview, movie.PosterPath, movie.GenreIds, movie.ID)

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

func (pg *PostgresMovieStore) DeleteMovieById(id int64) error {
	query := `
	DELETE FROM movies
	WHERE id = $1
	`

	result, err := pg.db.Exec(query, id)

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
