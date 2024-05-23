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
			http.Error(w, "You need to be logged in to view profiles", http.StatusUnauthorized)
			return
		}

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

		isOwnProfile := loggedInUserID == userID

		data := map[string]interface{}{
			"User":         user,
			"UserPosts":    userPosts,
			"LikedPosts":   likedPosts,
			"IsOwnProfile": isOwnProfile,
			"LoggedIn":     true,
		}

		utils.RenderTemplate(w, "profile/viewProfile.html", data)
	}
}
