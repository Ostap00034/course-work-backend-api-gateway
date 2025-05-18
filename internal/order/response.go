// internal/order/response.go
package order

import commonpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func errorResponse(msg string, errs map[string]string) Response {
	return Response{Success: false, Message: msg, Errors: errs}
}

func successResponse(msg string) Response {
	return Response{Success: true, Message: msg}
}

type OrderResponse struct {
	Response
	Order *commonpbv1.OrderData `json:"order,omitempty"`
}

type OrdersResponse struct {
	Response
	Orders []*commonpbv1.OrderData `json:"orders"`
}

type MyOrdersResponse struct {
	Response
	Orders []*commonpbv1.OrderData `json:"orders,omitempty"`
}
