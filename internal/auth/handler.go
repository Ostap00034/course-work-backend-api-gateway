package auth

import (
	"net/http"

	util "github.com/Ostap00034/course-work-backend-api-gateway/util"
	authv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

// loginRequest — тело запроса для /login.
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginHandler
// @Summary      Авторизация
// @Description  Логин по email и паролю, выставляет httpOnly cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      loginRequest  true  "Параметры авторизации"
// @Success      200      {object}  Response       "успех"
// @Failure      400      {object}  Response       "ошибка валидации"
// @Failure      401      {object}  Response       "неверные логин/пароль"
// @Failure      500      {object}  Response       "внутренняя ошибка"
// @Router       /auth/login [post]
func LoginHandler(client authv1.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					switch fe.Field() {
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
					case "Fio":
						if fe.Tag() == "required" {
							errs["fio"] = "фио обязательно"
						} else {
							errs["fio"] = "минимальная длина фио 4 символа"
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

		// RPC
		resp, err := client.Login(c, &authv1.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				c.JSON(http.StatusNotFound, Response{
					Success: false,
					Message: "неверный логин или пароль",
				})
			} else {
				c.JSON(http.StatusInternalServerError, Response{
					Success: false,
					Message: "внутренняя ошибка сервера",
				})
			}
			return
		}

		util.SetCookie(c, "token", resp.Token, resp.ExpiresAt)
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "авторизация прошла успешно",
		})
	}
}

// ValidateHandler
// @Summary      Проверка токена
// @Description  Проверяет валидность текущего httpOnly cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} Response "токен валидный"
// @Failure      401 {object} Response "токен невалидный или истек"
// @Router       /auth/validate [get]
func ValidateHandler(client authv1.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := util.GetCookie(c, "token")
		if !ok {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "токен обязателен",
			})
			return
		}
		resp, err := client.ValidateToken(c, &authv1.ValidateTokenRequest{Token: token})
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "токен невалидный или истек",
			})
			return
		}

		user := resp.User

		c.JSON(http.StatusOK, Response{
			Success:   true,
			Message:   "токен валидный",
			ExpiresAt: resp.ExpiresAt,
			User:      user,
			UserId:    resp.UserId,
		})
	}
}

// LogoutHandler
// @Summary      Выход
// @Description  Отозвать токен и очистить cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} Response "выход успешен"
// @Failure      500 {object} Response "выход не удался"
// @Router       /auth/logout [post]
func LogoutHandler(client authv1.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := util.GetCookie(c, "token")
		_, err := client.Revoke(c, &authv1.RevokeRequest{Token: token})
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Message: "выход не удался",
			})
			return
		}
		util.ClearCookie(c, "token")
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "выход успешен",
		})
	}
}

// RegisterHandlers вешает маршруты /auth.
func RegisterHandlers(r gin.IRouter, client authv1.AuthServiceClient) {
	r.POST("/login", LoginHandler(client))
	r.GET("/validate", ValidateHandler(client))
	r.POST("/logout", LogoutHandler(client))
}
