/*
server/main.go
Author: Bipin Kumar Ojha
Description: This is the main entry point of the server application. It initializes the configuration, logger, database, Firebase, and HTTP router, then starts the server on the configured port.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhilckenshi/notification/internal/database"
	routers "github.com/akhilckenshi/notification/internal/routes"
	"github.com/akhilckenshi/notification/pkg/logger"

	cfg "github.com/akhilckenshi/notification/pkg/settings"
)

func main() {
	// Initialize application configuration from environment variables and settings files
	config, err := initializeConfig()
	if err != nil {
		fmt.Println("Not able to get config files")
	}

	// Initialize the logger with settings from the configuration
	initializeLogger(config)
	logger.Log.Info("Logger Initialized")

	// Initialize the database
	if err := database.InitDatabase(config); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer func() {
		if err := database.CloseDatabase(); err != nil {
			logger.Log.Error(fmt.Sprintf("Error closing database: %v", err))
		}
	}()

	// Initialize the HTTP router with the registered routes
	router := routers.GetRouter()
	logger.Log.Info("Router Initialized")

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if config.App.WithSSL {
			certFile := "/etc/ssl/certs/cert.pem"
			keyFile := "/etc/ssl/private/key.pem"
			// Start the HTTPS server on the configured port, log and handle any errors
			if err := router.ListenTLS(fmt.Sprintf(":%s", config.AppPort), certFile, keyFile); err != nil {
				logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
			}
		} else {
			// Start the HTTP server on the configured port, log and handle any errors
			if err := router.Listen(fmt.Sprintf(":%s", config.AppPort)); err != nil {
				logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
			}
		}
	}()

	// Wait for a termination signal before gracefully closing services
	<-quit
	logger.Log.Info("Shutting down server...")

	logger.Log.Info("Server gracefully stopped.")
}

/*
initializeConfig loads the application configuration from environment variables
and other settings files using the InitConfig function from the settings package.
*/
func initializeConfig() (cfg.Configuration, error) {
	return cfg.InitConfig()
}

/*
initializeLogger sets up the logging system using the settings specified
in the configuration file (e.g., log file name, size, retention, and level).
*/
func initializeLogger(conf cfg.Configuration) {
	logger.InitLogger(
		conf.Logger.FileName,
		conf.Logger.FileSize,
		conf.Logger.MaxLogFile,
		conf.Logger.MaxRetention,
		conf.Logger.CompressLog,
		conf.Logger.Level,
	)
}
