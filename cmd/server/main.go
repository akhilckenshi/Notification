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
	"smartsme-notificationservice/internal/database"
	routers "smartsme-notificationservice/internal/routes"
	"smartsme-notificationservice/pkg/logger"
	"syscall"

	cfg "smartsme-notificationservice/pkg/settings"
)

func main() {
	// Initialize application configuration from environment variables and settings files
	initializeConfig()

	// Initialize the logger with settings from the configuration
	initializeLogger()
	logger.Log.Info("Logger Initialized")

	// Initialize the database
	if err := database.InitDatabase(cfg.Config.Database); err != nil {
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
		// Start the HTTP server on the configured port, log and handle any errors
		if err := router.Listen(fmt.Sprintf(":%d", cfg.Config.Server.Port)); err != nil {
			logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
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
func initializeConfig() {
	cfg.InitConfig()
}

/*
initializeLogger sets up the logging system using the settings specified
in the configuration file (e.g., log file name, size, retention, and level).
*/
func initializeLogger() {
	logger.InitLogger(
		cfg.Config.Logger.FileName,
		cfg.Config.Logger.FileSize,
		cfg.Config.Logger.MaxLogFile,
		cfg.Config.Logger.MaxRetention,
		cfg.Config.Logger.CompressLog,
		cfg.Config.Logger.Level,
	)
}


