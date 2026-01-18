package router

import (
	"net/http"

	chirps "github.com/lucasgjanot/go-http-server/internal/api/v1/chirps"
	chirpsid "github.com/lucasgjanot/go-http-server/internal/api/v1/chirps/chirpid"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/healthz"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/login"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/metrics"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/polka/webhooks"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/refresh"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/reset"
	"github.com/lucasgjanot/go-http-server/internal/api/v1/revoke"
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
	// admin
	mux.HandleFunc("GET /admin/metrics", metrics.GetHandler(cfg.Metrics))
	mux.HandleFunc("POST /admin/reset", reset.PostHandler(cfg.Metrics))
	// status
	mux.HandleFunc("GET /api/healthz", healthz.GetHandler())
	// chirps
	mux.HandleFunc("GET /api/chirps", chirps.GetHandler())
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirpsid.GetHandler())
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", chirpsid.DeleteHandler(cfg.JWTSecret))
	mux.HandleFunc("POST /api/chirps", chirps.PostHandler(cfg.JWTSecret))
	// users
	mux.HandleFunc("POST /api/users", users.PostHandler())
	mux.HandleFunc("PUT /api/users", users.PutHandler(cfg.JWTSecret))
	// login
	mux.HandleFunc("POST /api/login", login.PostHandler(cfg.JWTSecret))

	mux.HandleFunc("POST /api/refresh", refresh.PostHandler(cfg.JWTSecret))
	mux.HandleFunc("POST /api/revoke", revoke.PostHandler())

	mux.HandleFunc("POST /api/polka/webhooks", webhooks.PostHandler(cfg.PolkaAPIKey))
	
	return &Router{
		Handler: mux,
	}
}