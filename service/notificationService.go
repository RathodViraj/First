package service

import (
	"First/model"
	"First/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(r repository.NotificationRepository) *NotificationService {
	return &NotificationService{r}
}

func (s *NotificationService) SaveNotification(noti model.Notification) error {
	return s.repo.SaveNotification(noti)
}

func (s *NotificationService) GetNotificationsByUser(userID int) ([]model.Notification, error) {
	return s.repo.GetNotificationsByUser(userID)
}
