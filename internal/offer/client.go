package offer

import (
	offerpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/offer/v1"
	"google.golang.org/grpc"
)

func NewClient(cc *grpc.ClientConn) offerpbv1.OfferServiceClient {
	return offerpbv1.NewOfferServiceClient(cc)
}
