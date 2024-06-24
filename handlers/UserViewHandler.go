

package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db" 
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

		userIDStr := segments[2] 
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		loggedInUserID, err := GetUserIDFromSession(r)
		if err != nil {
			loggedInUserID = 0
		}

		options := map[string]bool{
			"userPosts":     true,
			"likedPosts":    true,
			"notifications": true,
		}

		
		data, err := utils.GetPageData(dbConn, loggedInUserID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := db.GetUserByID(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}
		userPosts, err := db.GetAllPosts(dbConn, 0, int64(userID), "likes")
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

		data.Title = user.Username + "'s Profile"
		data.User = user
		data.IsOwnProfile = isOwnProfile
		data.Posts = userPosts 
		data.LikedPosts = likedPosts

		data.LoggedIn = loggedInUserID > 0

		utils.RenderTemplate(w, "profile/profile.html", data)
	}
}
