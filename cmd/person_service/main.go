package main

import (
	"log/slog"
	"person_service/internal/config"
	"person_service/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.Setup(cfg.Environment)

	log.Info("Starting person_service", slog.String("environment", cfg.Environment))

	if cfg.Environment == config.EnvLocal || cfg.Environment == config.EnvDev {
		log.Debug("Debug messages are enabled")
	}

}
