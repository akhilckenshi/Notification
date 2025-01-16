/*
database/database.go
Author: Bipin Kumar Ojha
Description: This file defines the Database interface, which outlines the necessary methods for initializing and closing database connections. It also provides functionality to initialize the appropriate database based on the configuration and to close an active connection.
*/

package database

import (
	"fmt"

	"github.com/akhilckenshi/notification/pkg/logger"
	cfg "github.com/akhilckenshi/notification/pkg/settings"
)

// Database interface defines methods that any database implementation should provide
// Initialize: Establish a connection to the database
// Close: Close the established database connection
type Database interface {
	Initialize() error // Initializes the database connection
	Close() error      // Closes the database connection
	GetDBClient() any  // Method to get database client
	GetDBName() string // Method to get the database name
}

// dbHandler holds the active database connection object, which implements the Database interface
var dbHandler Database

// InitDatabase initializes the appropriate database connection based on the configuration
// It checks the database type from the config file and sets the dbHandler accordingly
func InitDatabase(conf cfg.Configuration) error {
	fmt.Printf("CONFIG: %+v\n", conf)
	// Switch case to determine the type of database being used, such as MongoDB or PostgreSQL
	switch conf.Database.DBType {
	case "mongo":
		dbHandler = &MongoDB{
			DBUri:    conf.DBURI,       // MongoDB URI
			DBName:   conf.DBName,      // MongoDB Database Name
			MaxLimit: conf.DBConnCount, // Maximum number of connections
		}
	default:
		// If the database type is unsupported, return an error
		return fmt.Errorf("unsupported database type: %s", conf.Database.DBType)
	}

	// Initialize the database connection and return any error encountered during initialization
	return dbHandler.Initialize()
}

// CloseDatabase gracefully closes the current database connection if one exists
func CloseDatabase() error {
	if dbHandler == nil {
		// Log a message indicating that there is no active database connection to close
		logger.Log.Debug("No active database connection to close")
		return nil
	}
	// Close the database connection and return any error encountered during closure
	return dbHandler.Close()
}

// GetDBClient retrieves the database client from the dbHandler
func GetDBClient() any {
	if dbHandler == nil {
		return nil // Handle the case where there is no active database connection
	}
	return dbHandler.GetDBClient()
}

// GetDBName retrieves the database name from the dbHandler
func GetDBName() string {
	if dbHandler == nil {
		return "" // Handle the case where there is no active database connection
	}
	return dbHandler.GetDBName()
}
