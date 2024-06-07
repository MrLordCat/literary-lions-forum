package db

import (
	"database/sql"
	"log"
)

func AddNotification(db *sql.DB, userID int, postID int, content string) error {
	query := `INSERT INTO notifications (user_id, post_id, content) VALUES (?, ?, ?)`
	_, err := db.Exec(query, userID, postID, content)
	if err != nil {
		log.Printf("Error adding notification: %v", err)
		return err
	}
	return nil
}

func GetUnreadNotifications(db *sql.DB, userID int) ([]Notification, error) {
	query := `
        SELECT id, user_id, content, is_read, created_at, post_id 
        FROM notifications 
        WHERE user_id = ? AND is_read = 0
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Content, &notification.IsRead, &notification.CreatedAt, &notification.PostID)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func GetUnreadNotificationsCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = false", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func MarkNotificationsAsRead(db *sql.DB, userID int) error {
	query := `UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0`
	_, err := db.Exec(query, userID)
	if err != nil {
		log.Printf("Error marking notifications as read: %v", err)
		return err
	}
	return nil
}
