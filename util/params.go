package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUIDParam читает путь c.Param(name), парсит его в uuid.UUID
// и при ошибке сразу пишет 400 и abort’ит запрос.
func ParseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "invalid id format",
					"errors": map[string]string{name: "must be a valid UUID"},
			})
			c.Abort()
			return uuid.Nil, false
	}
	return id, true
}
