// handlers/UserViewHandler.go

package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db" // Импортируйте ваш пакет db, проверьте правильность пути
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// UserViewHandler отображает профиль указанного пользователя.
func UserViewHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL or User ID", http.StatusBadRequest)
		return
	}

	userIDStr := segments[2] // предполагается, что URL в формате /user/{userID}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Замените эту функцию на ваш метод получения соединения с БД

	user, err := db.GetUserByID(dbConn, userID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	userPosts, err := db.GetUserPosts(dbConn, userID)
	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		return
	}

	likedPosts, err := db.GetLikedPosts(dbConn, userID)
	if err != nil {
		http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/template/viewProfile.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"User":       user,
		"UserPosts":  userPosts,
		"LikedPosts": likedPosts,
	})
	if err != nil {
		http.Error(w, "Failed to render profile page", http.StatusInternalServerError)
	}
}
