package main

import (
	pb "commons/api"
	"context"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *grpcHandler) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenResponse, error) {
	log.Printf("Validating token called with token: %s\n", token.Token)
	return &pb.TokenResponse{
		Expired:   false,
		CreatedAt: "",
		ExpiredAt: "",
	}, nil
}
