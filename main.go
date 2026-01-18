package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/database"
	"github.com/lucasgjanot/go-http-server/internal/router"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	if env == "production" {
		envFile = ".env"
	} else {
		envFile = ".env." + env
	}

	// Load env file ONLY if it exists
	if err := godotenv.Load(envFile); err != nil {
		if env != "production" {
			log.Fatalf("Error loading %s: %v", envFile, err)
		}
		// production: ignore missing file
	}

	log.Printf("Running in %s mode", env)

	const filepathRoot = "./app"
	const port = "8080"

	cfg := config.NewConfig()
	database.Init(cfg.Database.DbURL)
	srv := router.New(cfg)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv.Handler,
	}

	log.Printf("Serving on address: http://localhost:%s\n", port)
	log.Fatal(httpServer.ListenAndServe())
}

