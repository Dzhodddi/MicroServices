package main

import (
	pb "commons/api"
	"commons/database"
	redisDB "commons/redis"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	grpcHandler struct {
		redisService *redisDB.RedisService
		pb.UnimplementedAuthServiceServer
	}
)

func NewGrpcHandler(grpcServer *grpc.Server, redis *redis.Client) {
	handler := &grpcHandler{
		redisService: &redisDB.RedisService{Conn: redis},
	}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *grpcHandler) ValidateToken(ctx context.Context, tokenRequest *pb.TokenRequest) (*pb.TokenResponse, error) {
	user, err := h.redisService.ValidateUserToken(ctx, tokenRequest.Email)
	if err != nil {
		switch {
		case errors.Is(err, database.NotFound):
			return nil, status.Error(codes.NotFound, "not found")
		default:
			return nil, status.Error(codes.Internal, "server error")
		}

	}
	return &pb.TokenResponse{
		Expired: false,
		Email:   user.Email,
		Ttl:     user.TTL,
	}, nil
}
