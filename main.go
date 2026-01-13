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

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("No .env.%s file found", env)
	}

	const filepathRoot = "./app"
	const port = "8080"

	cfg := config.NewConfig()
	database.Init(cfg.Database.DbURL)
	srv := router.New(cfg)

	httpServer := &http.Server{
		Addr: ":" + port,
		Handler: srv.Handler,
	}
	

	log.Printf("Serving on address: http://localhost:%s\n", port)
	log.Fatal(httpServer.ListenAndServe())
}
