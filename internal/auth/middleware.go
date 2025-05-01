package auth

import (
    "github.com/gin-gonic/gin"
    "github.com/Ostap00034/course-work-backend-api-gateway/util/cookie"
)

// Middleware прокидывает токен из cookie в gRPC-metadata.
func Middleware() gin.HandlerFunc {
    return util.CookieToMetadata()
}
