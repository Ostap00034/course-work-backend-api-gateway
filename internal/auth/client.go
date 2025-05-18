package auth

import (
	authpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	"google.golang.org/grpc"
)

// NewClient возвращает gRPC-клиент AuthService.
func NewClient(cc *grpc.ClientConn) authpbv1.AuthServiceClient {
	return authpbv1.NewAuthServiceClient(cc)
}
