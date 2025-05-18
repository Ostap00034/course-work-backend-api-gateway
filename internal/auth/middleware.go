package auth

import (
	"net/http"

	"github.com/Ostap00034/course-work-backend-api-gateway/util"
	"github.com/Ostap00034/course-work-backend-auth-service/util/jwt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// Middleware прокидывает токен из cookie в gRPC-metadata.
func Middleware() gin.HandlerFunc {
    return util.CookieToMetadata()
}

// AdminOnly проверяет, что в метадате gRPC (или JWT-claim) есть роль admin.
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        // вытянуть токен из cookie → metadata (у вас уже есть CookieToMetadata)
        // здесь просто читаем claims из контекста
        md, ok := metadata.FromIncomingContext(c.Request.Context())
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "нет метадаты"})
            return
        }
        toks := md["authorization"]
        if len(toks) == 0 {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "нет токена"})
            return
        }
        claims, err := jwt.ParseToken(toks[0])
        if err != nil || claims.Role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступ запрещён"})
            return
        }
        // прокидываем в контекст пользователя, если нужно
        c.Next()
    }
}
