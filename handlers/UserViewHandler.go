// handlers/UserViewHandler.go

package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db" // Импортируйте ваш пакет db, проверьте правильность пути
	"literary-lions-forum/utils"
	"net/http"
	"strconv"
	"strings"
)

func UserViewHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		loggedInUserID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
			return
		}

		options := map[string]bool{
			"userPosts":  true,
			"likedPosts": true,
		}

		// Получаем данные для профиля другого пользователя
		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := db.GetUserByID(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		isOwnProfile := loggedInUserID == userID

		data.Title = user.Username + "'s Profile"
		data.User = user
		data.IsOwnProfile = isOwnProfile
		data.Posts = data.UserPosts // Добавляем посты пользователя

		utils.RenderTemplate(w, "profile/profile.html", data)
	}
}
