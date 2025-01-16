/*
repo/notification.go
Author: Akhil C
Description: Repository for managing notification data in MongoDB.
*/

package repo

import (
	"context"
	"fmt"

	"time"

	"github.com/akhilckenshi/Notification/pkg/logger"
	"github.com/akhilckenshi/notification/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Warehouse handles interactions with the notification collection
type Notification struct {
	db *mongo.Collection
}

// NewNotificationRepo initializes the notification with a MongoDB collection
func NewNotificationRepo(cl interface{}, dbName string) *Notification {
	if mongoClient, ok := cl.(*mongo.Client); ok {
		collectionName := models.Notification{}.TableName()
		collection := mongoClient.Database(dbName).Collection(collectionName)

		return &Notification{db: collection}
	} else {
		return nil
	}
}

// Store Notification information from the Message
func (repo *Notification) StoreNotificationInformation(notification models.Notification) error {
	notification.UpdatedAt = time.Now()
	_, err := repo.db.InsertOne(context.TODO(), notification)
	if err != nil {
		errStr := fmt.Sprintf("failed to store information: %v", err)
		logger.Log.Error(errStr)
		return fmt.Errorf(errStr)
	}
	return nil
}

// Listnotifications lists all notifications.
func (repo *Notification) ListNotifications(ctx context.Context, filter bson.M) ([]*models.Notification, error) {
	cursor, err := repo.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var notificaitons []*models.Notification
	if err := cursor.All(ctx, &notificaitons); err != nil {
		return nil, err
	}
	return notificaitons, nil
}
