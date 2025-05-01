package user

import (
    "google.golang.org/grpc"
    userv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
)

// NewClient возвращает gRPC-клиент UserService.
func NewClient(cc *grpc.ClientConn) userv1.UserServiceClient {
    return userv1.NewUserServiceClient(cc)
}
