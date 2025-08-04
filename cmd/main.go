package main

import (
	"fmt"
	"github.com/AramLab/api-gateway/internal/clients"
	"github.com/AramLab/api-gateway/internal/config"
	"github.com/AramLab/api-gateway/internal/handlers"
	customLogger "github.com/AramLab/api-gateway/internal/logger"
	"github.com/AramLab/api-gateway/internal/router"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal(fmt.Errorf("error loading .env file: %w", err))
		return
	}

	var cfg config.APIGatewayConfig

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(fmt.Errorf("error processing .env file: %w", err))
		return
	}

	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing logger: %w", err))
		return
	}

	grpcClient, err := clients.NewAuthClient(cfg.Server.GRPCPort)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing grpcClient: %w", err))
		return
	}

	authHandler := handlers.NewAuthHandler(grpcClient, logger)
	defer grpcClient.Close()

	taskHandler := handlers.NewTaskHandler(logger, "http://localhost"+cfg.Server.ListenAddr)

	app := router.NewRouters(&router.Routers{TaskHandler: *taskHandler, AuthHandler: *authHandler}, cfg.Server.Secret)

	go func() {
		logger.Infof("Starting server on %s", cfg.Server.ListenAddr)
		if err = app.Listen(cfg.Server.ListenAddr); err != nil {
			log.Fatal(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down server...")
}
