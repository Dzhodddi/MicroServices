package database

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

var (
	QueryTimeOut      = 5 * time.Second
	ErrDuplicateEmail = errors.New("duplicate email")
	NotFound          = errors.New("not found")
)

func New(addr string, maxOpenConnections, maxIdleConnections int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeOut)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
