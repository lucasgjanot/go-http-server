package middleware

import (
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/config"
)

func Metrics(m *config.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Inc()
			next.ServeHTTP(w,r)
		})
	}
}