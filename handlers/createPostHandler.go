package handlers

import (
	"database/sql" // Импортируем твой пакет db
	"fmt"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

func PostCreateFormHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение ID пользователя из сессии
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
			return
		}
		if userID == 0 {
			http.Error(w, "User not logged in", http.StatusForbidden)
			return
		}

		// Получение данных из формы
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDString := r.FormValue("category_id")
		categoryID, err := strconv.Atoi(categoryIDString)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Создание поста с учётом категории
		if err := db.CreatePost(dbConn, title, content, userID, categoryID); err != nil {
			http.Error(w, "Failed to save post: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Post created:", title, content, categoryID)
		// Перенаправление на главную страницу после успешного создания
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
