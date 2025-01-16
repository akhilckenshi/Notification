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

	"time"

	"github.com/akhilckenshi/notification/internal/models"
	"github.com/akhilckenshi/notification/internal/notifications"
	"github.com/akhilckenshi/notification/internal/repo"
	config "github.com/akhilckenshi/notification/pkg/settings"
	"github.com/akhilckenshi/notification/pkg/utils"

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
	consumer, err := sarama.NewConsumer([]string{config.Config.KafkaPort}, configs)
	if err != nil {
		fmt.Println("Error creating Kafka consumer:", err)
		return
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition(config.Config.KafkaTopic, 0, sarama.OffsetNewest)
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
				msg.From = config.Config.Email.Id
				notifications.SendEmail(msg.To, msg.Subject, msg.Message)
			case "whatsapp":
				fmt.Println("Processing WhatsApp notification")
				msg.From = config.Config.WhatsAppFromNumber
				notifications.SendWhatsAppMessage(msg.To, msg.Subject, msg.Message)
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
		Subject:        notifier.Subject,
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

// // NotificationService handles business logic for notification
// type NotificationService struct {
// 	repo *repo.Notification
// }

// // NewNotificationService creates a new instance of NotificationService
// func NewNotificationService(repo *repo.Notification) *NotificationService {
// 	return &NotificationService{repo: repo}
// }

// // Consumer
// func (s *NotificationService) MessageConsumer() {
// 	time.Sleep(10 * time.Second)

// 	configs := sarama.NewConfig()
// 	configs.Consumer.Return.Errors = true
// 	configs.Net.DialTimeout = 10 * time.Second
// 	configs.Net.ReadTimeout = 10 * time.Second
// 	configs.Net.WriteTimeout = 10 * time.Second

// 	kafkaBrokers := []string{config.Config.KafkaPort}

// 	// Create consumer with retry logic
// 	var consumer sarama.Consumer
// 	var err error

// 	for i := 0; i < 5; i++ {
// 		consumer, err = sarama.NewConsumer(kafkaBrokers, configs)
// 		if err == nil {
// 			break
// 		}
// 		log.Printf("Failed to create consumer: %v. Retrying in 5 seconds...", err)
// 		time.Sleep(5 * time.Second)
// 	}

// 	if err != nil {
// 		fmt.Println("Error creating Kafka consumer:", err)
// 		return
// 	}
// 	defer func() {
// 		if err := consumer.Close(); err != nil {
// 			log.Printf("Error closing consumer: %v", err)
// 		}
// 	}()

// 	partitionConsumer, err := consumer.ConsumePartition(config.Config.KafkaTopic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		fmt.Println("Error creating partition consumer:", err)
// 		return
// 	}
// 	defer func() {
// 		if err := partitionConsumer.Close(); err != nil {
// 			log.Printf("Error closing partition consumer: %v", err)
// 		}
// 	}()

// 	fmt.Println("Kafka consumer started. Waiting for messages...")

// 	for {
// 		select {
// 		case message := <-partitionConsumer.Messages():
// 			msg, err := s.UnmarshelChatMessage(message.Value)
// 			if err != nil {
// 				fmt.Println("Error unmarshalling message:", err)
// 				continue
// 			}

// 			// Determine message content based on priority
// 			// var messageContent string
// 			// switch msg.Priority {
// 			// case "high":
// 			// 	messageContent = "Do it fast - High priority"
// 			// case "medium":
// 			// 	messageContent = "Do it soon - Medium priority"
// 			// case "low":
// 			// 	messageContent = "Do it when possible - Low priority"
// 			// default:
// 			// 	messageContent = "Do it as per priority"
// 			// }
// 			// Send notification based on message type and priority
// 			switch msg.Type {
// 			case "email":
// 				fmt.Println("Processing Email notification")
// 				msg.From = config.Config.Email.Id
// 				notifications.SendEmail(msg.To, msg.Subject, msg.Message)
// 			case "whatsapp":
// 				fmt.Println("Processing WhatsApp notification")
// 				msg.From = config.Config.WhatsAppFromNumber
// 				notifications.SendWhatsAppMessage(msg.To, msg.Subject, msg.Message)
// 			default:
// 				fmt.Println("Unknown message type:", msg.Type)
// 			}
// 			fmt.Println("Received message:", msg)
// 			err = s.repo.StoreNotificationInformation(*msg)
// 			if err != nil {
// 				fmt.Println("Error storing message in repository:", err)
// 				continue
// 			}
// 		case err := <-partitionConsumer.Errors():
// 			fmt.Println("Kafka consumer error:", err)
// 		}
// 	}
// }

// // Unmarshal byte to Notification structure from Notifier
// func (n *NotificationService) UnmarshelChatMessage(data []byte) (*models.Notification, error) {
// 	var notifier models.Notifier // Create an instance of Notifier for unmarshalling
// 	err := json.Unmarshal(data, &notifier)
// 	if err != nil {
// 		fmt.Println("Error unmarshalling Notifier:", err)
// 		return nil, err
// 	}

// 	// Map Notifier fields to Notification, assigning ID to NotificationID
// 	notification := &models.Notification{
// 		NotificationID: notifier.ID,
// 		OrganizationID: notifier.OrganizationID,
// 		To:             notifier.To,
// 		From:           notifier.From,
// 		Type:           notifier.Type,
// 		Priority:       notifier.Priority,
// 		Subject:        notifier.Subject,
// 		Message:        notifier.Message,
// 		Status:         "Delivered",
// 		CreatedAt:      notifier.CreatedAt,
// 		UpdatedAt:      time.Now(),
// 	}

// 	return notification, nil
// }

// // GetNotification retrieves a list of all notifications from the repository and returns them.
// func (s *NotificationService) GetNotifications(ctx context.Context, key, orgId string) ([]*models.Notification, error) {
// 	filter, _ := utils.GetNotificationFilter(key, orgId)

// 	return s.repo.ListNotifications(ctx, filter)
// }
