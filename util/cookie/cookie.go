package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// SetCookie выставляет HTTPOnly cookie.
func SetCookie(c *gin.Context, name, value string, expiresAt int64) {
	maxAge := int(expiresAt - time.Now().Unix())
	c.SetCookie(name, value, maxAge, "/", "", false, true)
}

// GetCookie читает HTTPOnly cookie.
func GetCookie(c *gin.Context, name string) (string, bool) {
	v, err := c.Cookie(name)
	return v, err == nil
}

// ClearCookie удаляет cookie.
func ClearCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", false, true)
}

// CookieToMetadata прокидывает значение cookie "auth_token" в gRPC-metadata.
func CookieToMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, ok := GetCookie(c, "auth_token"); ok {
			md := metadata.Pairs("authorization", token)
			c.Request = c.Request.WithContext(metadata.NewOutgoingContext(c.Request.Context(), md))
		}
		c.Next()
	}
}
