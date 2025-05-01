package user

import (
	"net/http"

	pb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/gin-gonic/gin"
)

// RegisterHandlers вешает маршруты /users.
func RegisterHandlers(r gin.IRouter, client pb.UserServiceClient) {
	r.POST("", func(c *gin.Context) {
		var req pb.CreateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.CreateUser(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userId": resp.UserId})
	})
}
