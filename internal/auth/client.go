package auth

import (
	authv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	"google.golang.org/grpc"
)

// NewClient возвращает gRPC-клиент AuthService.
func NewClient(cc *grpc.ClientConn) authv1.AuthServiceClient {
	return authv1.NewAuthServiceClient(cc)
}
