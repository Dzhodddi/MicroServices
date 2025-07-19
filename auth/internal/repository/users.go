package repository

import (
	"commons/database"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UsersRepository struct {
	db *sql.DB
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &password
	p.hash = hash
	return nil
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

type UserDB struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

func (repo *UsersRepository) create(ctx context.Context, tx *sql.Tx, user *UserDB) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return database.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (repo *UsersRepository) createInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExpr time.Duration, userID int64) error {
	query := `INSERT INTO user_inventations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, invitationExpr, userID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UsersRepository) Register(ctx context.Context, user *UserDB, token string, invitationExpr time.Duration) error {
	return withTx(repo.db, ctx, func(tx *sql.Tx) error {
		if err := repo.create(ctx, tx, user); err != nil {
			return err
		}

		if err := repo.createInvitation(ctx, tx, token, invitationExpr, user.ID); err != nil {
			return err
		}

		return nil
	})
}
