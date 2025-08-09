package store

import "database/sql"

type Review struct {
	ID        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	MovieId   int64  `json:"movie_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type PostgresReviewStore struct {
	db *sql.DB
}

func NewPostgresReviewStore(db *sql.DB) *PostgresReviewStore {
	return &PostgresReviewStore{db: db}
}

type ReviewStore interface {
	CreateReview(*Review) error
	GetReviewById(id int64) (*Review, error)
	UpdateReview(*Review) error
	DeleteReview(id int64) error
}

func (pg *PostgresReviewStore) CreateReview(review *Review) error {
	query := `
	INSERT INTO reviews (user_id, movie_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`
	err := pg.db.QueryRow(query, review.UserId, review.MovieId, review.Content).Scan(&review.ID, &review.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresReviewStore) GetReviewById(id int64) (*Review, error) {
	query := `
	SELECT id, user_id, movie_id, content
	FROM reviews
	WHERE id = $1
	`
	review := &Review{}
	err := pg.db.QueryRow(query, id).Scan(&review.ID, &review.UserId, &review.MovieId, &review.Content)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (pg *PostgresReviewStore) UpdateReview(review *Review) error {
	query := `
	UPDATE reviews
	SET content = $1
	WHERE id = $2
	`
	result, err := pg.db.Exec(query, review.Content, review.ID)
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

func (pg *PostgresReviewStore) DeleteReview(id int64) error {
	query := `
	DELETE FROM reviews
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
