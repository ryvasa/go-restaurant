package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type AuthenticationMiddleware struct {
	tokenUtil *utils.TokenUtil
}

func NewAuthenticationMiddleware(tokenUtil *utils.TokenUtil) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{tokenUtil: tokenUtil}
}

func (m *AuthenticationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.WithField("authentication_errors", nil).Error("Error invalid request header")
			utils.HttpResponse(w, http.StatusUnauthorized, nil,
				utils.NewUnauthorizedError("No authorization header"))
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Log.WithField("authentication_errors", nil).Error("Error invalid token format")
			utils.HttpResponse(w, http.StatusUnauthorized, nil,
				utils.NewUnauthorizedError("Invalid authorization format"))
			return
		}

		claims, err := m.tokenUtil.ExtractClaims(tokenParts[1])
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid token")
			utils.HttpResponse(w, http.StatusUnauthorized, nil,
				utils.NewUnauthorizedError("Invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
