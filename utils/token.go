package utils

import (
	"fmt"
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

func (t *TokenUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(t.config.Secret.JwtSecretKey), nil
    })
}

func (t *TokenUtil) ExtractClaims(tokenString string) (jwt.MapClaims, error) {
    token, err := t.ValidateToken(tokenString)
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}
