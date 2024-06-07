package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"log"
	"net/http"
)

func NotificationHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		options := map[string]bool{
			"notifications": true,
		}

		data, err := utils.GetPageData(database, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Пометка уведомлений как прочитанных
		if err := db.MarkNotificationsAsRead(database, userID); err != nil {
			http.Error(w, "Failed to mark notifications as read: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Title = "Notifications"

		log.Printf("Notifications: %+v", data.Notifications) // Логируем полученные уведомления

		utils.RenderTemplate(w, "notifications.html", data)
	}
}
func MarkNotificationsAsReadHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := db.MarkNotificationsAsRead(dbConn, userID); err != nil {
			http.Error(w, "Failed to mark notifications as read: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
