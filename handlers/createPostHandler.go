package handlers

import (
	"database/sql" // Импортируем твой пакет db
	"literary-lions-forum/handlers/db"
	"net/http"
)

func CreatePostHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получение ID пользователя из сессии
	userID, err := GetUserIDFromSession(r)
	if err != nil {
		// Обрабатывай ошибку, например, отправляя ответ клиенту
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}
	if userID == 0 {
		http.Error(w, "User not logged in", http.StatusForbidden)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")

	if err := db.CreatePost(dbConn, title, content, userID); err != nil {
		http.Error(w, "Failed to save post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
