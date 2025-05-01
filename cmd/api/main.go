package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Ostap00034/course-work-backend-api-gateway/internal/auth"
	"github.com/Ostap00034/course-work-backend-api-gateway/internal/user"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// 1) Gin + middleware
	r := gin.Default()
	api := r.Group("api")
	api.Use(auth.Middleware())

	// 2) gRPC–сonnections
	authSvcAddr, exists := os.LookupEnv("AUTH_SERVICE_ADDR")
	if !exists {
		log.Fatal("not AUTH_SERVICE_ADDR in .env file")
	}
	authConn, err := grpc.NewClient(authSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial auth-service: %v", err)
	}

	userSvcAddr, exists := os.LookupEnv("USER_SERVICE_ADDR")
	if !exists {
		log.Fatal("not USER_SERVICE_ADDR in .env file")
	}
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial user-service: %v", err)
	}

	// 3) Клиенты
	authClient := auth.NewClient(authConn)
	userClient := user.NewClient(userConn)

	// 4) Роуты по фичам
	auth.RegisterHandlers(api.Group("/auth"), authClient)
	user.RegisterHandlers(api.Group("/users"), userClient)

	// 5) Запуск
	addr, exists := os.LookupEnv("GATEWAY_ADDR")
	if !exists {
		log.Fatal("not GATEWAY_ADDR in .env file")
	}
	log.Printf("API Gateway listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
