package router

import (
	"net/http"

	chirps "github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/healthz"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/metrics"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/reset"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/users"
	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/middleware"
)

type Router struct {
	Handler http.Handler
}

func New(cfg *config.Config) *Router{
	const filepathRoot = "./app"

	mux := http.NewServeMux()
	metricsMiddleware := middleware.Metrics(cfg.Metrics)

	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux.Handle(
		"/app/",
		http.StripPrefix("/app", metricsMiddleware(fileServer)),
	)

	mux.HandleFunc("GET /admin/metrics", metrics.GetHandler(cfg.Metrics))
	mux.HandleFunc("POST /admin/reset", reset.PostHandler(cfg.Metrics))
	mux.HandleFunc("GET /api/healthz", healthz.GetHandler())
	mux.HandleFunc("POST /api/chirps", chirps.PostHandler())
	mux.HandleFunc("POST /api/users", users.PostHandler())

	return &Router{
		Handler: mux,
	}
}