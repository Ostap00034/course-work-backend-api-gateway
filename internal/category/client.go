package category

import (
	categorypbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/category/v1"
	"google.golang.org/grpc"
)

// NewClient возвращает gRPC-клиент AuthService.
func NewClient(cc *grpc.ClientConn) categorypbv1.CategoryServiceClient {
	return categorypbv1.NewCategoryServiceClient(cc)
}
