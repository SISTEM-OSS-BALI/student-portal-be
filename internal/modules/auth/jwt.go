package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/username/gin-gorm-api/internal/schema"
)

const (
	envJWTSecret    = "JWT_SECRET"
	envJWTExpiresIn = "JWT_EXPIRES_IN"
	envJWTIssuer    = "JWT_ISSUER"
)

const (
	defaultIssuer = "student-portal"
)

var defaultExpiresIn = 24 * time.Hour

type Config struct {
	Secret    []byte
	ExpiresIn time.Duration
	Issuer    string
}

type Claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   schema.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func LoadConfigFromEnv() (Config, error) {
	secret := os.Getenv(envJWTSecret)
	if secret == "" {
		return Config{}, fmt.Errorf("%s is required", envJWTSecret)
	}

	expiresIn := defaultExpiresIn
	if raw := os.Getenv(envJWTExpiresIn); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf("invalid %s: %w", envJWTExpiresIn, err)
		}
		expiresIn = parsed
	}

	issuer := os.Getenv(envJWTIssuer)
	if issuer == "" {
		issuer = defaultIssuer
	}

	return Config{
		Secret:    []byte(secret),
		ExpiresIn: expiresIn,
		Issuer:    issuer,
	}, nil
}

func GenerateToken(user schema.User) (string, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
		Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.ExpiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.Secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
