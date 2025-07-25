package repository

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
}

type Storage struct {
	Users interface {
		Register(ctx context.Context, user *UserDB, token string, invitationExpr time.Duration) error
		Login(ctx context.Context, email, pass string) (error, *UserDB)
		Activate(ctx context.Context, token string) error
	}
}

func withTx(db *sql.DB, ctx context.Context, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UsersRepository{db},
	}
}
