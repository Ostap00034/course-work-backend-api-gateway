package order

import (
	orderpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/order/v1"
	"google.golang.org/grpc"
)

// NewClient возвращает gRPC-клиент AuthService.
func NewClient(cc *grpc.ClientConn) orderpbv1.OrderServiceClient {
	return orderpbv1.NewOrderServiceClient(cc)
}
