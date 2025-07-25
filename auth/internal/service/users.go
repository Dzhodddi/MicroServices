package service

import (
	"auth/auth"
	"auth/internal/repository"
	"auth/shared"
	"commons/broker"
	"commons/database"
	"commons/redis"
	"commons/shared_errors"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	jwtExpr = 3 * time.Hour
)

type UserService struct {
	store  *repository.Storage
	broker broker.BrokerService
	redis  redis.RedisService
	auth   auth.Authenticator
}

func (s *UserService) RegisterNewUser(ctx context.Context, payload shared.RegisterNewUser) (*shared.UserWithToken, error) {

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])
	user := repository.UserDB{
		Username: payload.Username,
		Email:    payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		return nil, shared_errors.ServerError
	}
	err := s.store.Users.Register(ctx, &user, hashedToken, time.Hour)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDuplicateEmail):
			return nil, shared_errors.ViolatePK
		default:
			return nil, shared_errors.ServerError
		}
	}
	msg := map[string]string{
		"name":    payload.Username,
		"token":   plainToken,
		"address": "http://localhost:3050/v1/users/activate",
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, shared_errors.ServerError
	}
	err = s.broker.Publish(ctx, bytes, "email")
	if err != nil {
		return nil, shared_errors.ServerError
	}

	return &shared.UserWithToken{
		Username: payload.Username,
		Email:    payload.Email,
		Token:    plainToken,
	}, nil

}

func (s *UserService) Login(ctx context.Context, payload shared.LoginUser) (*shared.User, string, error) {
	err, user := s.store.Users.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, database.NotFound):
			return nil, "", shared_errors.NotFoundError
		default:
			fmt.Print(err.Error())
			return nil, "", shared_errors.ServerError
		}
	}
	claims := jwt.MapClaims{
		"sub": payload.Email,
		"exp": time.Now().Add(jwtExpr).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": issuer,
		"aud": audience,
	}

	token, err := s.auth.GenerateToken(claims)
	if err != nil {
		fmt.Print(err.Error())
		return nil, "", shared_errors.ServerError
	}

	err = s.redis.SetUserToken(ctx, user.Email, token, jwtExpr)
	if err != nil {
		return nil, "", shared_errors.ServerError
	}
	return &shared.User{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		IsActive:  user.IsActive,
		RoleID:    user.RoleID,
	}, token, nil
}

func (s *UserService) Activate(ctx context.Context, token string) error {
	err := s.store.Users.Activate(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, database.NotFound):
			return shared_errors.NotFoundError
		default:
			fmt.Print(err)
			return shared_errors.ServerError
		}
	}

	return nil
}
