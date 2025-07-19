package main

import (
	"auth/internal/repository"
	"auth/internal/service"
	"commons"
	"commons/broker"
	commonDB "commons/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

var (
	grpcAddr = commons.EnvString("GRPC_ADDR", "localhost:2000")
	httpAddr = commons.EnvString("HTTP_ADDR", "localhost:3000")
	version  = "0.0.1"
	env      = commons.EnvString("ENV", "development")
)

//	@title			Auth microserver
//	@description	API for auth
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath		/v1
//
// @description	Auth microserver
var dbAddr = commons.EnvString("DB_ADDR", "")
var brokerAddr = commons.EnvString("BROKER_ADDR", "")

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// http handlers
		database, err := commonDB.New(dbAddr, 5, 5, "15m")
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		brokerConnection, err := broker.New(brokerAddr)
		if err != nil {
			log.Fatalf("failed to connect to broker: %v", err)
		}
		store := repository.NewStorage(database)

		services := service.NewService(&store, brokerConnection)

		srv := handler{service: services, server: echo.New()}
		srv.mount()

		srv.server.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: 5 * time.Second,
		}))
		srv.server.Use(middleware.Logger())
		srv.server.Use(middleware.Recover())
		srv.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

		log.Println("Listening http server on:", httpAddr)

		if err = srv.server.Start(httpAddr); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		//grpc server
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		defer func(listener net.Listener) {
			err = listener.Close()
			if err != nil {
				log.Fatalf("failed to close listener: %v", err)
			}
		}(listener)
		grpcServer := grpc.NewServer()
		NewGrpcHandler(grpcServer)

		log.Println("Listening grpc server on:", grpcAddr)
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	wg.Wait()
	log.Printf("All Auth service is stopped")
}
