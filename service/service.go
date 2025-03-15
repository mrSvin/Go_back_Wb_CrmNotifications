package service

import (
	"backWbCrmNotifications/models"
	"backWbCrmNotifications/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.repo.Create(notification)
}

func (s *NotificationService) GetAllNotifications() ([]models.Notification, error) {
	return s.repo.GetAll()
}
