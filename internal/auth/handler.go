// internal/auth/handler.go
package auth

import (
	"net/http"

	util "github.com/Ostap00034/course-work-backend-api-gateway/util/cookie"
	authv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTO для /login
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterHandlers(r gin.IRouter, client authv1.AuthServiceClient) {
	r.POST("/login", func(c *gin.Context) {
		var req loginRequest

		// 1) Валидация JSON + binding + validator.v10
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			// validator.ValidationErrors содержит детали
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					field := fe.Field()
					switch field {
					case "Email":
						if fe.Tag() == "required" {
							errs["email"] = "электронная почта обязательна"
						} else {
							errs["email"] = "неверный формат электронной почты"
						}
					case "Password":
						if fe.Tag() == "required" {
							errs["password"] = "пароль обязателен"
						} else {
							errs["password"] = "минимальная длина пароля 6 символов"
						}
					}
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "ошибка валидации",
				Errors:  errs,
			})
			return
		}

		// 2) RPC-вызов
		resp, err := client.Login(c, &authv1.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.Unauthenticated:
					c.JSON(http.StatusUnauthorized, Response{
						Success: false,
						Message: st.Message(),
					})
				default:
					c.JSON(http.StatusInternalServerError, Response{
						Success: false,
						Message: st.Message(),
					})
				}
			} else {
				c.JSON(http.StatusInternalServerError, Response{
					Success: false,
					Message: st.Message(),
				})
			}
			return
		}

		// 3) Успех
		util.SetCookie(c, "auth_token", resp.Token, resp.ExpiresAt)
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "авторизация прошла успешно",
		})
	})

	// остальные хендлеры (validate, logout) можно аналогично обернуть в Response
}
