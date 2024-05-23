package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
	"text/template"
)

func NotificationHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получите ID пользователя из сессии или другого источника
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		notifications, err := db.GetUnreadNotifications(database, userID)
		if err != nil {
			http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("web/template/notifications.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, notifications)
	}
}
