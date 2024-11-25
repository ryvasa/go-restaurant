package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type AuthorizationMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationMiddleware(enforcer *casbin.Enforcer) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{enforcer: enforcer}
}

func (m *AuthorizationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user role from context (set by auth middleware)
		claims, ok := r.Context().Value("user").(jwt.MapClaims)
		if !ok {
			logger.Log.WithField("user", r.Context().Value("user")).Error("Error invalid user claims")

			utils.HttpResponse(w, http.StatusUnauthorized, nil,
				utils.NewUnauthorizedError("Invalid token claims"))
			return
		}

		role := claims["role"].(string)
		path := r.URL.Path
		method := r.Method

		// Check permission
		allowed, err := m.enforcer.Enforce(role, path, method)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid permission")
			utils.HttpResponse(w, http.StatusInternalServerError, nil,
				utils.NewInternalError("Authorization check failed"))
			return
		}
		fmt.Println(role, path, method)
		if !allowed {
			logger.Log.WithError(err).Error("Error not allowed")
			utils.HttpResponse(w, http.StatusForbidden, nil,
				utils.NewUnauthorizedError("Insufficient permissions"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
