package main

import (
	pb "commons/api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

var (
	httpAddr        = os.Getenv("HTTP_ADDR")
	authServiceAddr = os.Getenv("AUTH_SERVICE_ADDR")
)

func main() {

	conn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	client := pb.NewAuthServiceClient(conn)
	handler := NewHandler(client)
	srv := echo.New()
	srv.Use(middleware.Recover())
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${status}] ${method} ${path} ${latency_human}\n",
	}))
	srv.POST("/api/token", handler.ValidateToken)
	log.Println("Starting server on: ", httpAddr)

	if err = http.ListenAndServe(httpAddr, srv); err != nil {
		log.Fatal(err)
	}
}
