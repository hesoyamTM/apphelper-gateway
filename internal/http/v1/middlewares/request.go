package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
)

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), logger.RequestID, middleware.GetReqID(r.Context()))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
