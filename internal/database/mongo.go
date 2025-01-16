/*
database/mongo.go
Author: Bipin Kumar Ojha
Description: This file defines the MongoDB struct and its associated methods to initialize and close a MongoDB connection. It implements the Database interface.
*/

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/akhilckenshi/notification/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct holds the MongoDB client and configuration parameters
type MongoDB struct {
	Client   *mongo.Client // MongoDB client instance
	DBUri    string        // URI of the MongoDB server
	DBName   string        // Name of the database
	MaxLimit int           // Maximum connection pool size
}

// Initialize sets up a MongoDB connection and assigns it to the MongoDB struct
// Implements the Initialize method from the Database interface
func (m *MongoDB) Initialize() error {
	// Create a context with a 10-second timeout to avoid long-running operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled after the operation completes

	// Set client options, including the MongoDB URI and maximum connection pool size
	clientOptions := options.Client().ApplyURI(m.DBUri).SetMaxPoolSize(uint64(m.MaxLimit))

	// Attempt to establish a connection to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err) // Return error if connection fails
	}

	// Check if the MongoDB server is reachable by pinging it
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err) // Return error if ping fails
	}

	// Assign the connected client to the MongoDB struct
	m.Client = client
	// Log a successful connection message
	logger.Log.Info(fmt.Sprintf("Connected to MongoDB: %s", m.DBUri))

	return nil // Return nil if initialization was successful
}

// Close gracefully closes the MongoDB connection
// Implements the Close method from the Database interface
func (m *MongoDB) Close() error {
	// Check if the MongoDB client is not nil before attempting to close
	if m.Client != nil {
		// Attempt to disconnect the MongoDB client
		if err := m.Client.Disconnect(context.TODO()); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB: %v", err) // Return error if disconnection fails
		}
		// Log a message indicating the MongoDB connection has been closed
		logger.Log.Info("MongoDB connection closed.")
	}
	return nil // Return nil if closure was successful
}

func (m *MongoDB) GetDBClient() any {
	// Return the MongoDB client
	return m.Client
}

func (m *MongoDB) GetDBName() string {
	return m.DBName
}
