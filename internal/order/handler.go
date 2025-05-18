// internal/order/handler.go
package order

import (
	"net/http"
	"strings"

	orderpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/order/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

func CreateOrderHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					errs[strings.ToLower(fe.Field())] = "неверное или отсутствует поле"
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			c.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		if _, err := uuid.Parse(req.ClientId); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат client_id", nil))
			return
		}

		if _, err := uuid.Parse(req.CategoryId); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат categories_ids", nil))
			return
		}

		resp, err := client.CreateOrder(c, &orderpbv1.CreateOrderRequest{
			Title:       req.Title,
			Description: req.Description,
			Price:       req.Price,
			Address:     req.Address,
			Longitude:   req.Longitude,
			Latitude:    req.Latitude,
			CategoryId:  req.CategoryId,
			ClientId:    req.ClientId,
		})
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

		c.JSON(http.StatusOK, OrderResponse{
			Response: Response{Success: true, Message: "успешно"},
			Order:    resp.Order,
		})
	}
}

func GetOrdersHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getOrdersRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("некорректные параметры запроса", nil))
			return
		}

		// валидация только, без присваивания в лишнюю переменную
		if req.ClientId != "" {
			if _, err := uuid.Parse(req.ClientId); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат client_id", nil))
				return
			}
		}
		if req.MasterId != "" {
			if _, err := uuid.Parse(req.MasterId); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат master_id", nil))
				return
			}
		}

		for _, id := range req.CategoriesIds {
			if _, err := uuid.Parse(id); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат categories_ids", nil))
				return
			}
		}

		if req.ClientId == "" {
			req.ClientId = uuid.Nil.String()
		}

		if req.MasterId == "" {
			req.MasterId = uuid.Nil.String()
		}

		if req.CategoriesIds == nil {
			req.CategoriesIds = []string{}
		}

		resp, err := client.GetOrders(c, &orderpbv1.GetOrdersRequest{
			CategoriesIds: req.CategoriesIds,
			Status:        req.Status,
			ClientId:      req.ClientId,
			MasterId:      req.MasterId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}

		c.JSON(http.StatusOK, OrdersResponse{
			Response: Response{Success: true, Message: "успешно"},
			Orders:   resp.Orders,
		})
	}
}

func GetOrderHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат id", nil))
			return
		}
		resp, err := client.GetOrderById(c, &orderpbv1.GetOrderByIdRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}
		c.JSON(http.StatusOK, OrderResponse{
			Response: Response{Success: true, Message: "успешно"},
			Order:    resp.Order,
		})
	}
}

func UpdateOrderHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат id", nil))
			return
		}

		var req updateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := make(map[string]string)
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					errs[strings.ToLower(fe.Field())] = "неверное значение"
				}
			} else {
				errs["body"] = "некорректный запрос"
			}
			c.JSON(http.StatusBadRequest, errorResponse("ошибка валидации", errs))
			return
		}

		if req.ClientId != "" {
			if _, err := uuid.Parse(req.ClientId); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат client_id", nil))
				return
			}
		}
		if req.MasterId != "" {
			if _, err := uuid.Parse(req.MasterId); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат master_id", nil))
				return
			}
		}
		if req.CategoryId != "" {
			if _, err := uuid.Parse(req.CategoryId); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат category_id", nil))
				return
			}
		}

		resp, err := client.UpdateOrder(c, &orderpbv1.UpdateOrderRequest{
			Id:          id,
			Title:       req.Title,
			Description: req.Description,
			Address:     req.Address,
			Longitude:   req.Longitude,
			Latitude:    req.Latitude,
			Status:      req.Status,
			Price:       req.Price,
			CategoryId:  req.CategoryId,
			ClientId:    req.ClientId,
			MasterId:    req.MasterId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}

		c.JSON(http.StatusOK, OrderResponse{
			Response: Response{Success: true, Message: "успешно"},
			Order:    resp.Order,
		})
	}
}

func DeleteOrderHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат id", nil))
			return
		}
		if _, err := client.DeleteOrder(c, &orderpbv1.DeleteOrderRequest{Id: id}); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func GetMyOrdersHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getMyOrdersRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("некорректные параметры", nil))
			return
		}
		if _, err := uuid.Parse(req.UserId); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат user_id", nil))
			return
		}
		for _, id := range req.CategoriesIds {
			if _, err := uuid.Parse(id); err != nil {
				c.JSON(http.StatusBadRequest, errorResponse("неверный формат categories_ids", nil))
				return
			}
		}

		resp, err := client.GetMyOrders(c, &orderpbv1.GetMyOrdersRequest{
			UserId:        req.UserId,
			Status:        req.Status,
			CategoriesIds: req.CategoriesIds,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}
		c.JSON(http.StatusOK, MyOrdersResponse{
			Response: Response{Success: true, Message: "успешно"},
			Orders:   resp.Orders,
		})
	}
}

func GetMyFinishedOrdersHandler(client orderpbv1.OrderServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getMyFinishedOrdersRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("некорректные параметры", nil))
			return
		}
		if _, err := uuid.Parse(req.UserId); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("неверный формат user_id", nil))
			return
		}

		resp, err := client.GetMyFinishedOrders(c, &orderpbv1.GetMyFinishedOrdersRequest{
			UserId: req.UserId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("внутренняя ошибка сервера", nil))
			return
		}
		c.JSON(http.StatusOK, MyOrdersResponse{
			Response: Response{Success: true, Message: "успешно"},
			Orders:   resp.Orders,
		})
	}
}

func RegisterHandlers(r gin.IRouter, client orderpbv1.OrderServiceClient) {
	r.POST("/", CreateOrderHandler(client))
	r.GET("/", GetOrdersHandler(client))
	r.GET("/:id", GetOrderHandler(client))
	r.PUT("/:id", UpdateOrderHandler(client))
	r.DELETE("/:id", DeleteOrderHandler(client))
	r.GET("/my", GetMyOrdersHandler(client))
	r.GET("/my/finished", GetMyFinishedOrdersHandler(client))
}
