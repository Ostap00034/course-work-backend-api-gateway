package user

import commonpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type ProfileResponse struct {
	Response
	User *commonpb.UserData `json:"user,omitempty"`
}

type GetUsersResponse struct {
	Response
	Users []*commonpb.UserData `json:"users,omitempty"`
}

func errorResponse(msg string, errs map[string]string) Response {
	return Response{Success: false, Message: msg, Errors: errs}
}

func successResponse(msg string) Response {
	return Response{Success: true, Message: msg}
}
