package main

import (
	"log"
	"net/http"

	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/router"
)

func main() {
	const filepathRoot = "./app"
	const port = "8080"

	cfg := config.NewConfig()
	srv := router.New(cfg)

	httpServer := &http.Server{
		Addr: ":" + port,
		Handler: srv.Handler,
	}
	

	log.Printf("Serving on address: http://localhost:%s\n", port)
	log.Fatal(httpServer.ListenAndServe())
}