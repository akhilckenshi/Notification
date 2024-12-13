/*
models/notification.go
Author: Akhil C
Description: This file contains the document model and related structures for document information management.
*/

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`      // Unique identifier for the Notification
	NotificationID primitive.ObjectID `json:"notification_id" bson:"notification_id"` // ID of the notification
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"` // ID of the organization
	To             string             `json:"to" bson:"to"`                           // Recipient of the notification
	From           string             `json:"from" bson:"from"`                       // Sender of the notification
	Type           string             `json:"type" bson:"type"`                       // Type of the notification message
	Priority       string             `json:"priority" bson:"priority"`               // Priority level of the notification
	Subject        string             `json:"subject" bson:"subject"`                 // Subject of the notification message
	Message        string             `json:"message" bson:"message"`                 // Content of the notification message
	Status         string             `json:"status" bson:"status"`                   // Current status of the notification (e.g., sent, pending)
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`           // Timestamp of when the notification was created
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`           // Timestamp of when the notification was last updated
}

func (N Notification) TableName() string {
	return "notifiers" // Returns the collection name as 'Notifications'
}

type Notifier struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`      // Unique identifier for the Notification
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"` // ID of the organization
	To             string             `json:"to" bson:"to"`                           // Reciver of the notification
	From           string             `json:"from" bson:"from"`                       // Sender of the notification
	Type           string             `json:"type" bson:"type"`                       // Type of the notification message
	Priority       string             `json:"priority" bson:"priority"`               // Priority level of the notification
	Subject        string             `json:"subject" bson:"subject"`                 // Subject of the notification message
	Message        string             `json:"message" bson:"message"`                 // Content of the notification message
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`             // Timestamp of when the Business Type was created
}
