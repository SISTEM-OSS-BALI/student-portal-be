package router

import (
	neturl "net/url"
	"os"
	pathpkg "path"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	allowedOrigins := resolveAllowedOrigins()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return isOriginAllowed(origin, allowedOrigins)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func resolveAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	if raw == "" {
		return []string{"http://localhost:3000"}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		value := normalizeOriginPattern(part)
		if value == "" {
			continue
		}
		origins = append(origins, value)
	}

	if len(origins) == 0 {
		return []string{"http://localhost:3000"}
	}

	return origins
}

func normalizeOriginPattern(value string) string {
	return strings.TrimRight(strings.TrimSpace(value), "/")
}

func isOriginAllowed(origin string, patterns []string) bool {
	origin = normalizeOriginPattern(origin)
	if origin == "" {
		return true
	}

	originURL, err := neturl.Parse(origin)
	if err != nil {
		return false
	}

	for _, pattern := range patterns {
		if pattern == origin {
			return true
		}

		if !strings.Contains(pattern, "*") {
			continue
		}

		if matched, _ := pathpkg.Match(pattern, origin); matched {
			return true
		}

		patternURL, err := neturl.Parse(pattern)
		if err != nil {
			continue
		}

		if patternURL.Scheme != "" && patternURL.Scheme != originURL.Scheme {
			continue
		}

		hostPattern := patternURL.Hostname()
		if hostPattern == "" {
			hostPattern = strings.TrimPrefix(pattern, patternURL.Scheme+"://")
		}

		if matched, _ := pathpkg.Match(hostPattern, originURL.Hostname()); matched {
			return true
		}
	}

	return false
}
