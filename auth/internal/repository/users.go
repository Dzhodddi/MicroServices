package repository

import (
	"commons/database"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
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
	RoleID    int64    `json:"role_id"`
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

func (repo *UsersRepository) Login(ctx context.Context, email, pass string) (error, *UserDB) {
	user := &UserDB{}
	return withTx(repo.db, ctx, func(tx *sql.Tx) error {
		query := `SELECT id, username, email, is_active, created_at, role_id, password FROM users WHERE email = $1 and is_active = true`

		newCtx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
		defer cancel()
		var hashedPassword string
		err := tx.QueryRowContext(newCtx, query, email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.IsActive,
			&user.CreatedAt,
			&user.RoleID,
			&hashedPassword,
		)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return database.NotFound
			default:
				return err
			}
		}
		user.Password.hash = []byte(hashedPassword)
		if bcrypt.CompareHashAndPassword(user.Password.hash, []byte(pass)) != nil {
			return database.NotFound
		}
		return nil
	}), user

}

func (repo *UsersRepository) Activate(ctx context.Context, token string) error {
	return withTx(repo.db, ctx, func(tx *sql.Tx) error {

		user, err := repo.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		if err = repo.update(ctx, tx, user); err != nil {
			return err
		}
		if err = repo.deleteInvite(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (repo *UsersRepository) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*UserDB, error) {
	query := `
			SELECT u.id, u.username, u.email, u.created_at, u.is_active
			FROM users u
			JOIN user_inventations ui ON u.id = ui.user_id
			WHERE ui.token = $1 AND ui.expiry >= $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	user := &UserDB{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.NotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (repo *UsersRepository) update(ctx context.Context, tx *sql.Tx, user *UserDB) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UsersRepository) deleteInvite(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_inventations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UsersRepository) createInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExpr time.Duration, userID int64) error {
	query := `INSERT INTO user_inventations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, database.QueryTimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExpr))
	if err != nil {
		return err
	}
	return nil
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
