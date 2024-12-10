/*
service/notification.go
Author: Akhil C
Description: Service to manage notification data in MongoDB.
*/

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"smartsme-notificationservice/internal/models"
	"smartsme-notificationservice/internal/notifications"
	"smartsme-notificationservice/internal/repo"
	config "smartsme-notificationservice/pkg/settings"
	"smartsme-notificationservice/pkg/utils"
	"time"

	"github.com/IBM/sarama"
)

// NotificationService handles business logic for notification
type NotificationService struct {
	repo *repo.Notification
}

// NewNotificationService creates a new instance of NotificationService
func NewNotificationService(repo *repo.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}

// Consumer
func (s *NotificationService) MessageConsumer() {
	configs := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{config.Config.Kafka.KafkaPort}, configs)
	if err != nil {
		fmt.Println("Error creating Kafka consumer:", err)
		return
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition(config.Config.Kafka.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
		return
	}
	defer partitionConsumer.Close()
	fmt.Println("Kafka consumer started")
	for {
		select {
		case message := <-partitionConsumer.Messages():
			msg, err := s.UnmarshelChatMessage(message.Value)
			if err != nil {
				fmt.Println("Error unmarshalling message:", err)
				continue
			}

			// Determine message content based on priority
			// var messageContent string
			// switch msg.Priority {
			// case "high":
			// 	messageContent = "Do it fast - High priority"
			// case "medium":
			// 	messageContent = "Do it soon - Medium priority"
			// case "low":
			// 	messageContent = "Do it when possible - Low priority"
			// default:
			// 	messageContent = "Do it as per priority"
			// }
			// Send notification based on message type and priority
			switch msg.Type {
			case "email":
				fmt.Println("Processing Email notification")
				msg.From = config.Config.EMail.Id
				notifications.SendEmail(msg.To, "Reorder the quantity", msg.Message)
			case "whatsapp":
				fmt.Println("Processing WhatsApp notification")
				msg.From = config.Config.WhatsappConfig.Number
				notifications.SendWhatsAppMessage(msg.To, msg.Message)
			default:
				fmt.Println("Unknown message type:", msg.Type)
			}
			fmt.Println("Received message:", msg)
			err = s.repo.StoreNotificationInformation(*msg)
			if err != nil {
				fmt.Println("Error storing message in repository:", err)
				continue
			}
		case err := <-partitionConsumer.Errors():
			fmt.Println("Kafka consumer error:", err)
		}
	}
}

// Unmarshal byte to Notification structure from Notifier
func (n *NotificationService) UnmarshelChatMessage(data []byte) (*models.Notification, error) {
	var notifier models.Notifier // Create an instance of Notifier for unmarshalling
	err := json.Unmarshal(data, &notifier)
	if err != nil {
		fmt.Println("Error unmarshalling Notifier:", err)
		return nil, err
	}

	// Map Notifier fields to Notification, assigning ID to NotificationID
	notification := &models.Notification{
		NotificationID: notifier.ID,
		OrganizationID: notifier.OrganizationID,
		To:             notifier.To,
		From:           notifier.From,
		Type:           notifier.Type,
		Priority:       notifier.Priority,
		Message:        notifier.Message,
		Status:         "Delivered",
		CreatedAt:      notifier.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	return notification, nil
}

// GetNotification retrieves a list of all notifications from the repository and returns them.
func (s *NotificationService) GetNotifications(ctx context.Context, key, orgId string) ([]*models.Notification, error) {
	filter, _ := utils.GetNotificationFilter(key, orgId)

	return s.repo.ListNotifications(ctx, filter)
}
