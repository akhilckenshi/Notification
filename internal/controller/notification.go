/*
controller/notification.go
Author: Akhil C
Description: Controller to manage Notification data in MongoDB.
*/
package controller

import (
	"fmt"

	"github.com/akhilckenshi/notification/internal/service"
	"github.com/gofiber/fiber/v2"
)

// NotificationController defines HTTP handlers for Notifications.
type NotificationController struct {
	service *service.NotificationService
}

func NewNotificationController(service *service.NotificationService) *NotificationController {
	return &NotificationController{service: service}
}

func (c *NotificationController) ReadAllNotifications(ctx *fiber.Ctx) error {
	searchKey := ctx.Query("key")
	orgId := ctx.Query("orgID")

	// if searchKey == "" {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "key is required",
	// 	})
	// }

	if orgId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "organization ID is required",
		})
	}

	notificaitons, err := c.service.GetNotifications(ctx.Context(), searchKey, orgId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(notificaitons)
}

func Read(data string) string {
	return fmt.Sprintf("Hello module, I am %s", data)
}
