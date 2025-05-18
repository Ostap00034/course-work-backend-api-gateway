package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	categorypbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/category/v1"
	commonpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"
)

var validate = validator.New()

func CreateCategoryHandler(client categorypbv1.CategoryServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createCategoryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					switch fe.Field() {
					case "Name":
						if fe.Tag() == "required" {
							errs["name"] = "название обязательно"
						}
					case "Description":
						if fe.Tag() == "required" {
							errs["description"] = "описание обязательно"
						}
					}
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			ctx.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		resp, err := client.CreateCategory(ctx, &categorypbv1.CreateCategoryRequest{
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(st.Message(), nil))
			} else {
				msg := "внутренняя ошибка сервера"
				if st != nil {
					msg = st.Message()
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(msg, nil))
			}

			return
		}

		ctx.JSON(http.StatusOK, CategoryResponse{
			Response: Response{Success: true, Message: "успешно"},
			Category: resp.Category,
		})
	}
}

func GetCategoriesHandler(client categorypbv1.CategoryServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := client.GetCategories(ctx, &categorypbv1.GetCategoriesRequest{})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				ctx.JSON(http.StatusInternalServerError, errorResponse(st.Message(), nil))
			} else {
				ctx.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			}
			return
		}
		ctx.JSON(http.StatusOK, CategoriesResponse{
			Response:   Response{Success: true, Message: "успешно"},
			Categories: res.Categories,
		})
	}
}

func GetCategoryHandler(client categorypbv1.CategoryServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		categoryID, err := uuid.Parse(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse("неправильный формат id категории", nil))
			return
		}

		resp, err := client.GetCategoryById(ctx, &categorypbv1.GetCategoryByIdRequest{
			Id: categoryID.String(),
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				ctx.JSON(http.StatusInternalServerError, errorResponse(st.Message(), nil))
			} else {
				ctx.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			}
			return
		}

		ctx.JSON(http.StatusOK, CategoryResponse{
			Response: Response{Success: true, Message: "успешно"},
			Category: resp.Category,
		})
	}
}

func UpdateCategoryHandler(client categorypbv1.CategoryServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		categoryID, err := uuid.Parse(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse("неправильный формат id категории", nil))
			return
		}

		var req updateCategoryRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					switch fe.Field() {
					case "Name":
						errs["name"] = "название обязательно"
					case "Description":
						errs["description"] = "описание обязательно"
					}
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			ctx.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		categoryData := &commonpbv1.CategoryData{}
		if req.Name != "" {
			categoryData.Name = req.Name
		}
		if req.Description != "" {
			categoryData.Description = req.Description
		}

		resp, err := client.UpdateCategory(ctx, &categorypbv1.UpdateCategoryRequest{
			Id:       categoryID.String(),
			Category: categoryData,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				ctx.JSON(http.StatusInternalServerError, errorResponse(st.Message(), nil))
			} else {
				ctx.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			}
			return
		}

		ctx.JSON(http.StatusOK, CategoryResponse{
			Response: Response{Success: true, Message: "успешно"},
			Category: resp.Category,
		})
	}
}

func RegisterHandlers(r gin.IRouter, client categorypbv1.CategoryServiceClient) {
	r.POST("/", CreateCategoryHandler(client))
	r.GET("/", GetCategoriesHandler(client))
	r.GET("/:id", GetCategoryHandler(client))
	r.PUT("/:id", UpdateCategoryHandler(client))
	r.DELETE("/:id", GetCategoryHandler(client))
}
