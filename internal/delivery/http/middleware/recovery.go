package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
	"github.com/sirupsen/logrus"
)

type RecoveryConfig struct {
	EnableStackTrace bool
	LogError         bool
}

func RecoveryMiddleware(config *RecoveryConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = &RecoveryConfig{
			EnableStackTrace: true,
			LogError:         true,
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Buat error message yang informatif
					errorMsg := fmt.Sprintf("panic recovered: %v", err)
					logger.Log.WithFields(logrus.Fields{
						"error":     err,
						"path":      r.URL.Path,
						"method":    r.Method,
						"client_ip": r.RemoteAddr,
					}).Error(errorMsg)

					// Log error dengan stack trace jika diaktifkan
					if config.LogError {
						if config.EnableStackTrace {
							logger.Log.WithFields(logrus.Fields{
								"error":      err,
								"stacktrace": string(debug.Stack()),
								"path":       r.URL.Path,
								"method":     r.Method,
								"client_ip":  r.RemoteAddr,
							}).Error("Panic recovered in request handler")
						} else {
							logger.Log.WithFields(logrus.Fields{
								"error":     err,
								"path":      r.URL.Path,
								"method":    r.Method,
								"client_ip": r.RemoteAddr,
							}).Error("Panic recovered in request handler")
						}
					}

					// Kirim response error ke client
					utils.HttpResponse(w, http.StatusInternalServerError, nil,
						utils.AppError{
							HttpStatus: http.StatusInternalServerError,
							Code:       "INTERNAL_SERVER_ERROR",
							Message:    "An unexpected error occurred",
						})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
