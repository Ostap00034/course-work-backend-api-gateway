// @title        Course Work API
// @version      1.0
// @description  API Gateway для микросервисов
// @host         localhost:8080
// @BasePath     /api
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/Ostap00034/course-work-backend-api-gateway/cmd/api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Ostap00034/course-work-backend-api-gateway/internal/auth"
	"github.com/Ostap00034/course-work-backend-api-gateway/internal/category"
	"github.com/Ostap00034/course-work-backend-api-gateway/internal/offer"
	"github.com/Ostap00034/course-work-backend-api-gateway/internal/order"
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

	categorySvcAddr, exists := os.LookupEnv("CATEGORY_SERVICE_ADDR")
	if !exists {
		log.Fatal("not CATEGORY_SERVICE_ADDR in .env file")
	}
	categoryConn, err := grpc.NewClient(categorySvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial category-service: %v", err)
	}

	orderSvcAddr, exists := os.LookupEnv("ORDER_SERVICE_ADDR")
	if !exists {
		log.Fatal("not ORDER_SERVICE_ADDR in .env file")
	}
	orderConn, err := grpc.NewClient(orderSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial order-service: %v", err)
	}

	offerSvcAddr, exists := os.LookupEnv("OFFER_SERVICE_ADDR")
	if !exists {
		log.Fatal("not OFFER_SERVICE_ADDR in .env file")
	}
	offerConn, err := grpc.NewClient(offerSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial offer-service: %v", err)
	}

	// 3) Клиенты
	authClient := auth.NewClient(authConn)
	userClient := user.NewClient(userConn)
	categoryClient := category.NewClient(categoryConn)
	orderClient := order.NewClient(orderConn)
	offerClient := offer.NewClient(offerConn)

	// 4) Роуты по фичам
	auth.RegisterHandlers(api.Group("/auth"), authClient)
	user.RegisterHandlers(api.Group("/users"), userClient)
	category.RegisterHandlers(api.Group("/categories"), categoryClient)
	order.RegisterHandlers(api.Group("/orders"), orderClient)

	// 5) Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket для offer
	hub := offer.NewHub()
	api.GET("/ws/offers", offer.OfferWsHandler(hub, offerClient, authClient))

	// 6) Запуск
	addr, exists := os.LookupEnv("GATEWAY_ADDR")
	if !exists {
		log.Fatal("not GATEWAY_ADDR in .env file")
	}
	log.Printf("API Gateway listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
