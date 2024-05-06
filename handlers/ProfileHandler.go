package handlers

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func UserProfileHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromSession(r) // Получение ID пользователя из сессии
	if err != nil {
		http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(dbConn, userID) // Загрузка данных пользователя
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	userPosts, err := db.GetUserPosts(dbConn, userID) // Загрузка постов пользователя
	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		return
	}

	likedPosts, err := db.GetLikedPosts(dbConn, userID) // Загрузка лайкнутых постов
	if err != nil {
		http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
		return
	}

	// Рендеринг страницы профиля с полученными данными
	tmpl := template.Must(template.ParseFiles("web/template/profile.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"User":       user,
		"UserPosts":  userPosts,
		"LikedPosts": likedPosts,
	})
	if err != nil {
		http.Error(w, "Failed to render profile page", http.StatusInternalServerError)
	}
}
