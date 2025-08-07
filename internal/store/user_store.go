package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)

	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE username = $1
	`

	err := pg.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.CreatedAt)
	//returning now rows is not an error, it just means there is no data for the query
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
				UPDATE users
				SET username = $1, email = $2
				WHERE id = $3
			`
	result, err := pg.db.Exec(query, user.Username, user.Email, user.ID)

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

func (pg *PostgresUserStore) GetUserToken(scope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.created_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 and t.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := pg.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
