/*
routers.go
Author: Akhil C
Description: This file defines the main router configuration for the Fiber web application.
*/

package routers

import (
	"github.com/akhilckenshi/notification/internal/controller"
	"github.com/akhilckenshi/notification/internal/database"
	"github.com/akhilckenshi/notification/internal/repo"
	"github.com/akhilckenshi/notification/internal/service"
	"github.com/akhilckenshi/notification/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRouter initializes and returns the main Fiber application with configured routes.
func GetRouter() *fiber.App {
	app := fiber.New() // Initialize a new Fiber app

	// Create an API group for versioning or common routes..
	api := app.Group("/api")

	// Setup API version 1 (v1) routes..
	getV1ApiList(api)

	return app // Return the configured Fiber app..
}

// getV1ApiList sets up the version 1 (v1) API routes under the /api/v1 group.
func getV1ApiList(api fiber.Router) {
	// Group v1 routes under /api/v1..
	v1 := api.Group("/v1")

	// Initialize repositories and services..
	// Get the database client and database name..
	var notificationRepo *repo.Notification

	dbClient := database.GetDBClient()
	dbName := database.GetDBName()

	// Depending on the type of the database client, use MongoDB or PostgreSQL.
	if mongoClient, ok := dbClient.(*mongo.Client); ok {
		// MongoDB client.
		notificationRepo = repo.NewNotificationRepo(mongoClient, dbName)

	} else {
		// No database client available, log an error.
		logger.Log.Error("No DB Client available")
	}

	// Setup routes for Notification APIs.
	getNotificationApi(v1, notificationRepo)
}

// getNotificationApi sets up the Notification-related routes under /Account.
func getNotificationApi(v fiber.Router, notificationRepo *repo.Notification) {
	// Initialize Notification service and controller.
	notificationService := service.NewNotificationService(notificationRepo)
	notificationController := controller.NewNotificationController(notificationService)

	// Concurrently execute the messageConsumer
	go notificationService.MessageConsumer()

	// Define routes for Notification-related actions (Get).
	doc := v.Group("/notification")

	// Notification routes
	doc.Get("/", notificationController.ReadAllNotifications) // Route to retrieve all notifications from the system.
}
