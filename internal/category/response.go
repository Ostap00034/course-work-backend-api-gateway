package category

import commonpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type CategoryResponse struct {
	Response
	Category *commonpbv1.CategoryData `json:"category,omitempty"`
}

type CategoriesResponse struct {
	Response
	Categories []*commonpbv1.CategoryData `json:"categories"`
}

func errorResponse(msg string, errs map[string]string) Response {
	return Response{Success: false, Message: msg, Errors: errs}
}

func successResponse(msg string) Response {
	return Response{Success: true, Message: msg}
}
