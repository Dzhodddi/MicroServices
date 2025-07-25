package service

import (
	"auth/auth"
	"auth/internal/repository"
	"auth/shared"
	"commons/broker"
	redisDB "commons/redis"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"os"
)

const (
	issuer   = "Auth microserver"
	audience = "All microservers"
)

var secret = os.Getenv("JWT_SECRET")

type Service struct {
	IUserService interface {
		RegisterNewUser(ctx context.Context, payload shared.RegisterNewUser) (*shared.UserWithToken, error)
		Login(ctx context.Context, payload shared.LoginUser) (*shared.User, string, error)
		Activate(ctx context.Context, token string) error
	}
}

func NewService(store *repository.Storage, brokerConnection *amqp091.Connection, redisConnection *redis.Client) *Service {
	return &Service{
		IUserService: &UserService{
			store: store,
			broker: broker.BrokerService{
				Conn: brokerConnection,
			},
			redis: redisDB.RedisService{
				Conn: redisConnection,
			},
			auth: auth.NewJWTAuth(secret, audience, issuer),
		},
	}
}
