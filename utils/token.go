package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-restaurant/pkg/config"
)

type TokenUtil struct {
	config *config.Config
}

func NewTokenUtil(cfg *config.Config) *TokenUtil {
	return &TokenUtil{
		config: cfg,
	}
}

func (t *TokenUtil) GenerateToken(id, role string) (string, error) {
	claims := jwt.MapClaims{
		"iss":  "go-restaurant-api",
		"sub":  id,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.Secret.JwtSecretKey))
}
