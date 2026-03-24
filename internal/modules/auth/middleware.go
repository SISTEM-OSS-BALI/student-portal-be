package auth

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		serviceToken := strings.TrimSpace(c.GetHeader("X-Service-Token"))
		if isValidServiceToken(serviceToken) {
			c.Set("auth_service", true)
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "invalid Authorization header", nil)
				c.Abort()
				return
			}
			token = strings.TrimSpace(parts[1])
			if isValidServiceToken(token) {
				c.Set("auth_service", true)
				c.Next()
				return
			}
		} else {
			cookieToken, ok := readAuthCookie(c)
			if !ok {
				httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth token", nil)
				c.Abort()
				return
			}
			token = cookieToken
		}

		claims, err := ParseToken(token)
		if err != nil {
			httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "invalid or expired token", nil)
			c.Abort()
			return
		}

		c.Set("auth", claims)
		c.Next()
	}
}

func isValidServiceToken(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}

	expected := strings.TrimSpace(os.Getenv("SERVICE_TOKEN"))
	if expected == "" {
		expected = strings.TrimSpace(os.Getenv("AI_SERVICE_TOKEN"))
	}
	if expected == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}
