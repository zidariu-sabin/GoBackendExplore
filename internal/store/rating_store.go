package store

import "database/sql"

type Rating struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	MovieID int64 `json:"movie_id"`
	Score   int64 `json:"score"`
}

type PostgresRatingStore struct {
	db *sql.DB
}

func NewPostgresRatingStore(db *sql.DB) *PostgresRatingStore {
	return &PostgresRatingStore{db: db}
}

type RatingStore interface {
	CreateRating(rating *Rating) error
	GetMovieRatingScore(movieID int64) (float64, int, error)
	UpdateRating(rating *Rating) error
	DeleteRating(userId int64, movieID int64) error
}

func (pg *PostgresRatingStore) CreateRating(rating *Rating) error {
	query := `INSERT INTO ratings (user_id, movie_id, score) VALUES ($1, $2, $3) RETURNING id`
	err := pg.db.QueryRow(query, rating.UserID, rating.MovieID, rating.Score).Scan(&rating.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresRatingStore) GetMovieRatingScore(movieId int64) (float64, int, error) {
	query := `SELECT AVG(score)::NUMERIC(10,1), COUNT(*) FROM ratings WHERE movie_id = $1`

	var avgScore float64
	var ratingNo int
	err := pg.db.QueryRow(query, movieId).Scan(&avgScore, &ratingNo)

	if err != nil {
		return 0, 0, err
	}

	return avgScore, ratingNo, nil
}

func (pg PostgresRatingStore) UpdateRating(rating *Rating) error {
	query := `UPDATE ratings SET score = $1 WHERE id = $2`
	_, err := pg.db.Exec(query, rating.Score, rating.ID)

	return err
}

func (pg PostgresRatingStore) DeleteRating(userID int64, movieID int64) error {
	query := `DELETE FROM ratings WHERE user_id = $1 AND movie_id = $2`
	_, err := pg.db.Exec(query, userID, movieID)

	return err
}
