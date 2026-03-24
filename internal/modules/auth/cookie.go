package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const accessTokenCookieName = "sp_access_token"

func setAuthCookie(c *gin.Context, token string) {
	cfg, err := LoadConfigFromEnv()
	maxAge := 0
	if err == nil && cfg.ExpiresIn > 0 {
		maxAge = int(cfg.ExpiresIn.Seconds())
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		accessTokenCookieName,
		token,
		maxAge,
		"/",
		"",
		isSecureRequest(c),
		true,
	)
}

func readAuthCookie(c *gin.Context) (string, bool) {
	token, err := c.Cookie(accessTokenCookieName)
	if err != nil {
		return "", false
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}

func isSecureRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
