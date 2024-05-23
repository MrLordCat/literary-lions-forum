package server

import (
	"database/sql"
	"fmt"

	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
)

func MainPage(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := handlers.GetUserIDFromSession(r)
		loggedIn := err == nil && userID != 0

		// Получаем посты
		posts, err := db.GetAllPosts(dbConn, 0, 0)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Получаем категории
		categories, err := db.GetAllCategories(dbConn)
		if err != nil {
			http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
			return
		}

		// Получаем количество непрочитанных уведомлений
		var unreadNotifications int
		if loggedIn {
			unreadNotifications, err = db.GetUnreadNotificationsCount(dbConn, userID)
			if err != nil {
				//	http.Error(w, "Failed to fetch unread notifications", http.StatusInternalServerError)
				fmt.Println("Failed to fetch unread notifications")
			}
		}

		// Подготовка данных для шаблона
		data := struct {
			Title               string
			LoggedIn            bool
			UnreadNotifications int
			Categories          []db.Category
			Posts               []db.Post
			UserID              int
		}{
			Title:               "Welcome to the Book Forum",
			LoggedIn:            loggedIn,
			UnreadNotifications: unreadNotifications,
			Categories:          categories,
			Posts:               posts,
			UserID:              userID,
		}

		// Рендеринг шаблона
		utils.RenderTemplate(w, "home.html", data)
	}
}
