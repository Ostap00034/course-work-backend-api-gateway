package user

import (
	"net/http"

	commonpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"
	userv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

// DTOs

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Fio      string `json:"fio" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin master client"`
}

// updateUserRequest — все поля опциональны
type updateUserRequest struct {
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
	Fio   *string `json:"fio,omitempty" binding:"omitempty,min=4"`
	Role  *string `json:"role,omitempty" binding:"omitempty,oneof=admin master client"`
}

// Handlers

func CreateUserHandler(client userv1.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
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
					case "Fio":
						if fe.Tag() == "required" {
							errs["fio"] = "фио обязательно"
						} else {
							errs["fio"] = "минимальная длина фио 4 символа"
						}
					case "Password":
						if fe.Tag() == "required" {
							errs["password"] = "пароль обязателен"
						} else {
							errs["password"] = "минимальная длина пароля 6 символов"
						}
					case "Role":
						if fe.Tag() == "required" {
							errs["role"] = "роль обязателен"
						} else {
							errs["role"] = "неправильный формат роли"
						}
					}

				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			c.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		_, err := client.CreateUser(c, &userv1.CreateUserRequest{
			Email:    req.Email,
			Fio:      req.Fio,
			Role:     req.Role,
			Password: req.Password,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				c.JSON(http.StatusConflict, errorResponse(st.Message(), nil))
			} else {
				msg := "внутренняя ошибка сервера"
				if st != nil {
					msg = st.Message()
				}
				c.JSON(http.StatusInternalServerError, errorResponse(msg, nil))
			}
			return
		}

		c.JSON(http.StatusOK, successResponse("пользователь успешно создан"))
	}
}

func GetProfileHandler(client userv1.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неправильный формат id пользователя", nil))
			return
		}

		res, err := client.GetUserById(c, &userv1.GetUserByIdRequest{UserId: userID.String()})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				c.JSON(http.StatusNotFound, errorResponse(st.Message(), nil))
			} else {
				msg := "внутренняя ошибка сервера"
				if st != nil {
					msg = st.Message()
				}
				c.JSON(http.StatusInternalServerError, errorResponse(msg, nil))
			}
			return
		}

		c.JSON(http.StatusOK, ProfileResponse{
			Response: Response{Success: true, Message: "успешно"},
			User:     res.User,
		})
	}
}

func ChangeProfileHandler(client userv1.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неправильный формат id пользователя", nil))
			return
		}

		var req updateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					switch fe.Field() {
					case "Email":
						errs["email"] = "неверный формат электронной почты"
					case "Fio":
						errs["fio"] = "минимум 4 символа"
					case "Role":
						errs["role"] = "должно быть admin, master или client"
					}
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			c.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		// Формируем protobuf-структуру с изменениями
		userData := &commonpb.UserData{}
		if req.Email != nil {
			userData.Email = *req.Email
		}
		if req.Fio != nil {
			userData.Fio = *req.Fio
		}
		if req.Role != nil {
			userData.Role = *req.Role
		}

		res, err := client.ChangeUser(c, &userv1.ChangeUserRequest{
			UserId: userID.String(),
			User:   userData,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.NotFound:
					c.JSON(http.StatusNotFound, errorResponse(st.Message(), nil))
				case codes.InvalidArgument:
					c.JSON(http.StatusBadRequest, errorResponse(st.Message(), nil))
				default:
					c.JSON(http.StatusInternalServerError, errorResponse(st.Message(), nil))
				}
			} else {
				c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			}
			return
		}

		c.JSON(http.StatusOK, ProfileResponse{
			Response: Response{Success: true, Message: "успешно"},
			User:     res.User,
		})
	}
}

func GetUsersHandler(client userv1.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := client.GetUsers(c, &userv1.GetUsersRequest{})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				c.JSON(http.StatusInternalServerError, errorResponse(st.Message(), nil))
			} else {
				c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			}
			return
		}
		c.JSON(http.StatusOK, GetUsersResponse{
			Response: Response{Success: true, Message: "успешно"},
			Users:    res.Users,
		})
	}
}

func RegisterHandlers(r gin.IRouter, client userv1.UserServiceClient) {
	r.POST("/create", CreateUserHandler(client))
	r.GET("/profile/:id", GetProfileHandler(client))
	r.PUT("/profile/:id", ChangeProfileHandler(client))
	r.GET("/users", GetUsersHandler(client))
}
