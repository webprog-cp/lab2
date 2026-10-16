package main

import (
	"log"

	"cp_lab1/internal/app"
	"cp_lab1/internal/config"
	"cp_lab1/internal/router"

	_ "cp_lab1/docs"
)

// @title Cats shelter management API
// @version 1.0
// @description REST API service for cats shelter management.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	cfg := config.Load()
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r, err := router.New(cfg, a)
	if err != nil {
		log.Fatal(err)
	}

	host := cfg.Host + ":" + cfg.Port

	if err := r.Run(host); err != nil {
		log.Fatal(err)
	}
}
