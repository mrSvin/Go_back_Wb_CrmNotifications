package repository

import (
	"backWbCrmNotifications/models"
	"database/sql"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, title, url, message, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	row := r.db.QueryRow(query, notification.UserID, notification.Title, notification.URL, notification.Message, "В обработке")
	return row.Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) GetAll() ([]models.Notification, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, url, message, created_at, status FROM notifications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.URL, &n.Message, &n.CreatedAt, &n.Status)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
