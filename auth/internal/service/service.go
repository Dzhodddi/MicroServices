package service

import (
	"auth/internal/repository"
	"auth/shared"
	"commons"
	"commons/broker"
	"commons/database"
	"commons/shared_errors"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type Service struct {
	broker      broker.BrokerService
	store       *repository.Storage
	UserService interface {
	}
}

func NewService(store *repository.Storage, brokerConnection *amqp091.Connection) *Service {
	return &Service{
		broker: broker.BrokerService{
			Conn: brokerConnection,
		},
		store:       store,
		UserService: &UsersService{},
	}
}

func (s *Service) RegisterNewUser(c echo.Context) (*shared.User, error) {
	var payload shared.RegisterNewUser
	if err := c.Bind(&payload); err != nil {
		return nil, shared_errors.BadRequestPayload
	}

	if err := commons.Validate.Struct(&payload); err != nil {
		return nil, shared_errors.ValidationError
	}

	user := repository.UserDB{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return nil, shared_errors.ServerError
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := s.store.Users.Register(c.Request().Context(), &user, hashedToken, time.Hour*24)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDuplicateEmail):
			return nil, shared_errors.ViolatePK
		default:
			return nil, shared_errors.ServerError
		}
	}
	msg := map[string]string{
		"name":    user.Username,
		"token":   plainToken,
		"address": "http://localhost:3050/v1/users/activate",
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, shared_errors.ServerError
	}
	err = s.broker.Publish(c.Request().Context(), bytes, "email")
	if err != nil {
		return nil, shared_errors.ServerError
	}

	return &shared.User{
		Username: user.Username,
		Email:    user.Email,
		Token:    plainToken,
	}, nil

}
