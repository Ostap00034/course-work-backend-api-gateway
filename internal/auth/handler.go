// internal/auth/handler.go
package auth

import (
	"fmt"
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
		util.SetCookie(c, "token", resp.Token, resp.ExpiresAt)
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "авторизация прошла успешно",
		})
	})

	r.GET("/validate", func(c *gin.Context) {
		token, ok := util.GetCookie(c, "token")
		if !ok {
			c.JSON(http.StatusUnauthorized, errorResponse("токен обязателен", nil))
			return
		}
		resp, err := client.ValidateToken(c, &authv1.ValidateTokenRequest{Token: token})
		fmt.Println(resp)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse("токен невалидный или истек", nil))
			return
		}
		c.JSON(http.StatusOK, ValidateTokenResponse{
			Success:   true,
			Message:   "токен валидный",
			ExpiresAt: resp.ExpiresAt,
			UserId:    resp.UserId,
			Errors:    nil,
		})
	})

	r.POST("/logout", func(c *gin.Context) {
		token, _ := util.GetCookie(c, "token")
		_, err := client.Revoke(c, &authv1.RevokeRequest{Token: token})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("выход не удался", nil))
			return
		}
		util.ClearCookie(c, "token")
		c.JSON(http.StatusOK, successResponse("выход успешен"))
	})
}
