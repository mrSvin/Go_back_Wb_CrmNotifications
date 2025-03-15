package main

import (
	"backWbCrmNotifications/handlers"
	"backWbCrmNotifications/repository"
	"backWbCrmNotifications/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *sql.DB
	var err error

	// Повторные попытки подключения к базе данных
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to connect to the database: %v\n", i+1, err)
			time.Sleep(2 * time.Second) // Ждём 2 секунды перед следующей попыткой
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d: Failed to ping the database: %v\n", i+1, err)
			time.Sleep(2 * time.Second) // Ждём 2 секунды перед следующей попыткой
			continue
		}

		log.Println("Connected to the database")
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after %d attempts: %v", maxAttempts, err)
	}

	// Создаём таблицу, если она ещё не существует
	query := `
        CREATE TABLE IF NOT EXISTS notifications (
            id SERIAL PRIMARY KEY,
            user_id INT NOT NULL,
            title VARCHAR(255) NOT NULL,
            url VARCHAR(255) NOT NULL,
            message TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            status VARCHAR(50) DEFAULT 'В обработке'
        );
    `
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Table 'notifications' created or already exists")
	return db, nil
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	notificationRepo := repository.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			notificationHandler.CreateNotification(w, r)
		case http.MethodGet:
			notificationHandler.GetAllNotifications(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
