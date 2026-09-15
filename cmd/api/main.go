package main

import (
	_ "todo/docs"
	deliveryhttp "todo/internal/delivery/http"
	"todo/internal/infrastructure/container"
	"todo/internal/infrastructure/database"
	"todo/internal/infrastructure/server"
	"todo/pkg/config"
	"todo/pkg/logger"
)

// @title           Todo API
// @version         1.0
// @description     Simple Go Clean-Architecture Todo API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.AppEnv)

	db, err := database.NewPostgresConnection(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	c := container.NewContainer(db, log)

	router := deliveryhttp.SetupRouter(log, c)

	srv := server.New(router, cfg.ServerPort, log)
	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
