package user

import (
	"fmt"
	"net/http"

	pb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterHandlers вешает маршруты /users.
func RegisterHandlers(r gin.IRouter, client pb.UserServiceClient) {
	// Create User
	r.POST("", func(c *gin.Context) {
		var req createUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
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
			})
			return
		}
		resp, err := client.CreateUser(c, &pb.CreateUserRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.AlreadyExists:
					c.JSON(http.StatusBadRequest, Response{
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

		fmt.Println(resp)

		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "пользователь успешно создан",
		})
	})
}
