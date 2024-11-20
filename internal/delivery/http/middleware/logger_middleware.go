package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Get route pattern
        var match mux.RouteMatch
        route := ""
        if mux.CurrentRoute(r).Match(r, &match) {
            route = match.Route.GetName()
        }

        // Process request
        next.ServeHTTP(w, r)

        // Log after request is processed
        logger.Log.WithFields(logrus.Fields{
            "method":     r.Method,
            "path":      r.RequestURI,
            "route":     route,
            "duration":  time.Since(start),
            "ip":        r.RemoteAddr,
            "user_agent": r.UserAgent(),
        }).Info("HTTP request")
    })
}
