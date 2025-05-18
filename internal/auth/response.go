// internal/auth/response.go
package auth

import commonpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"

// Response — единый формат ответа.
type Response struct {
	Success   bool               `json:"success"`
	Message   string             `json:"message"`
	UserId    string             `json:"userId,omitempty"`
	User      *commonpb.UserData `json:"user,omitempty"`
	ExpiresAt int64              `json:"expiresAt,omitempty"`
	Errors    map[string]string  `json:"errors,omitempty"`
}
