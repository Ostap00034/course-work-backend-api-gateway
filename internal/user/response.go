package user

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
