package main

import (
	"commons"
	pb "commons/api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	httpAddr        = commons.EnvString("HTTP_ADDR", ":3000")
	authServiceAddr = "localhost:2000"
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
	srv.Use(middleware.Logger())
	srv.Use(middleware.Recover())
	srv.POST("/api/token/:token", handler.ValidateToken)
	log.Println("Starting server on: ", httpAddr)

	if err = http.ListenAndServe(httpAddr, srv); err != nil {
		log.Fatal(err)
	}
}
