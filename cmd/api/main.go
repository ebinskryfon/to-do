package main

import (
	deliveryhttp "todo/internal/delivery/http"
	"todo/internal/delivery/http/handlers"
	"todo/internal/infrastructure/server"
	"todo/pkg/config"
	"todo/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.AppEnv)

	// No database wiring yet — only the health check exists so far.
	healthHandler := handlers.NewHealthHandler()

	router := deliveryhttp.SetupRouter(log, healthHandler)

	srv := server.New(router, cfg.ServerPort, log)
	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
