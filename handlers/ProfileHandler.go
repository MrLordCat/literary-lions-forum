package handlers

import (
	"database/sql"
	"literary-lions-forum/utils"
	"net/http"
)

func UserProfileHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
			return
		}

		options := map[string]bool{
			"karma":         true,
			"userPosts":     true,
			"likedPosts":    true,
			"notifications": true,
			"isProfile":     true,
		}

		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.IsOwnProfile = true
		data.IsProfile = true // Установка флага IsProfile
		data.Posts = data.UserPosts

		utils.RenderTemplate(w, "profile/profile.html", data)
	}
}
