package handler

import (
	"First/model"
	"First/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{s}
}

func (h *NotificationHandler) SaveNotification(ctx *gin.Context) {
	notification := model.Notification{}
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.SaveNotification(notification); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to save notification"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Notification saved successfully"})
}

func (h *NotificationHandler) GetNotificationsByUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	if userIDStr == "" {
		ctx.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "User ID must be an integer"})
		return
	}

	notifications, err := h.service.GetNotificationsByUser(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	ctx.JSON(200, notifications)
}
