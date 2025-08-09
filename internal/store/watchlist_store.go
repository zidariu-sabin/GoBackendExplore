package store

import (
	"database/sql"
)

type WatchlistEntry struct {
	UserID  int64 `json:"user_id"`
	MovieID int64 `json:"movie_id"`
}

type PostgresWatchlistStore struct {
	db *sql.DB
}

func NewPostgresWatchlistStore(db *sql.DB) *PostgresWatchlistStore {
	return &PostgresWatchlistStore{
		db: db,
	}
}

type WatchlistStore interface {
	AddToWatchlist(userID int64, movieID int64) error
	RemoveFromWatchlist(userID int64, movieID int64) error
	GetWatchlist(userID int64) ([]Movie, error)
}

func (pg PostgresWatchlistStore) AddToWatchlist(userID int64, movieID int64) error {
	query := `
		INSERT into user_watchlist 
		(user_id, movie_id)
		VALUES ($1, $2)
	`

	_, err := pg.db.Exec(query, userID, movieID)

	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresWatchlistStore) RemoveFromWatchlist(userID int64, movieID int64) error {
	query := `
	DELETE FROM user_watchlist
	WHERE user_id = $1 AND movie_id = $2
	`

	_, err := pg.db.Exec(query, userID, movieID)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresWatchlistStore) GetWatchlist(userId int64) ([]Movie, error) {
	//TODO : add genre details
	query := `
		SELECT u.movie_id, m.Title, m.release_date, m.poster_path
		FROM user_watchlist u
		INNER JOIN movies m ON u.movie_id = m.id
		WHERE user_id = $1
		`

	var movies []Movie

	rows, err := pg.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var movie_entry Movie
		err = rows.Scan(&movie_entry.ID, &movie_entry.Title, &movie_entry.ReleaseDate, &movie_entry.PosterPath)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie_entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return movies, nil

}
