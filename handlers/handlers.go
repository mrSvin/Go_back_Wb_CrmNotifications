package handlers

import (
	"backWbCrmNotifications/models"
	"backWbCrmNotifications/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Логируем тело запроса для отладки
	log.Printf("Request body: %s", string(body))

	var notifications []models.Notification

	err = json.Unmarshal(body, &notifications)
	if err != nil {
		// Если не удалось декодировать как массив, пробуем декодировать как одиночное уведомление
		var singleNotification models.Notification
		err = json.Unmarshal(body, &singleNotification)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		notifications = []models.Notification{singleNotification}
	}

	// Создаём каждое уведомление в базе данных
	for _, notification := range notifications {
		// Устанавливаем значения по умолчанию
		notification.CreatedAt = time.Now() // Текущее время
		notification.Status = "В обработке" // Статус по умолчанию

		// Создаём уведомление в базе данных
		err := h.service.CreateNotification(&notification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := h.service.GetAllNotifications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
