/*
utils.go
Author: Akhil C
Description: This utility package provides common functions used across the application,
such as environment loading, random number generation, authentication utilities,
and JWT token generation.
*/

package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ObjectIDConversionError  = "Error converting string to ObjectID: %v"
	InvalidInputErrorMessage = "Invalid input data"
)

/*
LoadEnv loads environment variables from a .env file located in the specified directory.
It uses the godotenv package to read the .env file and set the environment variables.
If there is an error loading the .env file, the function logs the error and exits the program.
Parameters:
- currentConfigDirectory: The directory path where the .env file is located.
*/
func LoadEnv(currentConfigDirectory string) {
	// Construct the full path to the .env file.
	envPath := currentConfigDirectory + "/../../env/.env"
	fmt.Printf("Loading .env file from path: %s\n", envPath)

	// Attempt to load the .env file.
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file from path %s: %v\n", envPath, err)
	}
}

/*
GetEnv retrieves the value of the environment variable named by the key.
It returns the value of the specified environment variable, or an empty string if the key is not found.
Parameters:
- key: The name of the environment variable to retrieve.
Returns:
- The value of the environment variable as a string.
*/
func GetEnv(key string) string {
	// Fetch and return the value of the specified environment variable.
	return os.Getenv(key)
}

func GetNotificationFilter(key, organizationID string) (primitive.M, error) {
	// Convert organizationID from string to ObjectID
	orgObjID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %v", err)
	}

	// Base filter with organization ID as a mandatory field
	filter := bson.M{
		"organization_id": orgObjID, // Ensure organizationId is a required match
	}

	// Prepare $or clause for dynamic search on various fields
	orConditions := []bson.M{
		{"type": bson.M{"$regex": key, "$options": "i"}},
		{"priority": bson.M{"$regex": key, "$options": "i"}},
		{"from": bson.M{"$regex": key, "$options": "i"}},
		{"to": bson.M{"$regex": key, "$options": "i"}},
	}

	// Check if `key` is a valid ObjectID to add exact match conditions
	if objID, err := primitive.ObjectIDFromHex(key); err == nil {
		orConditions = append(orConditions, bson.M{"_id": objID})
		orConditions = append(orConditions, bson.M{"notification_id": objID})
	}

	// Add the $or clause to the filter as a sub-condition of $and
	filter["$and"] = []bson.M{
		{"organization_id": orgObjID},
		{"$or": orConditions},
	}

	return filter, nil
}
