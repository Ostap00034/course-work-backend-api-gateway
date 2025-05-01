// internal/auth/response.go
package auth

// Response — единый формат ответа.
type Response struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Errors  map[string]string `json:"errors,omitempty"`
}

// errorResponse упрощённо создаёт неуспешный ответ.
func errorResponse(msg string, errs map[string]string) Response {
    return Response{Success: false, Message: msg, Errors: errs}
}

// successResponse — для успешных случаев.
func successResponse(msg string) Response {
    return Response{Success: true, Message: msg}
}
