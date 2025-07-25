package main

import (
	"auth/docs"
	"auth/internal/repository"
	"auth/internal/service"
	"commons/broker"
	commonDB "commons/database"
	"commons/redis"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	grpcAddr   = ""
	httpAddr   = ""
	version    = "0.0.1"
	env        = ""
	dbAddr     = ""
	brokerAddr = ""
	apiUrl     = ""
	redisAddr  = ""
)

func init() {
	_ = godotenv.Load()
	grpcAddr = os.Getenv("GRPC_ADDR")
	httpAddr = os.Getenv("HTTP_ADDR")
	env = os.Getenv("ENV")
	dbAddr = os.Getenv("DB_ADDR")
	brokerAddr = os.Getenv("BROKER_ADDR")
	apiUrl = os.Getenv("API_URL")
	redisAddr = os.Getenv("REDIS_ADDR")
}

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
func main() {
	grpcAddr = os.Getenv("GRPC_ADDR")
	httpAddr = os.Getenv("HTTP_ADDR")
	version = "0.0.1"
	env = os.Getenv("ENV")
	dbAddr = os.Getenv("DB_ADDR")
	brokerAddr = os.Getenv("BROKER_ADDR")

	redisConnection := redis.NewRedisClient(redisAddr, "", 0)

	database, err := commonDB.New(dbAddr, 5, 5, "15m")
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	brokerConnection, err := broker.New(brokerAddr)
	if err != nil {
		log.Fatalf("failed to connect to broker: %v", err)
	}
	store := repository.NewStorage(database)

	services := service.NewService(&store, brokerConnection, redisConnection)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// http handlers

		srv := handler{service: services, server: echo.New()}
		docs.SwaggerInfo.Version = version
		docs.SwaggerInfo.Host = fmt.Sprint(apiUrl, httpAddr)
		docs.SwaggerInfo.BasePath = "/v1"
		srv.mount()

		srv.server.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: 5 * time.Second,
		}))

		srv.server.Use(middleware.Recover())
		srv.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
		srv.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${time_rfc3339} [${status}] ${method} ${path} ${latency_human}\n",
		}))
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
		NewGrpcHandler(grpcServer, redisConnection)

		log.Println("Listening grpc server on:", grpcAddr)
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	wg.Wait()
	log.Printf("All Auth service is stopped")
}
