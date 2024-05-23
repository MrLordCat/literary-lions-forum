package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
)

func UserProfileHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		karma, err := db.CalculateUserKarma(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to calculate karma", http.StatusInternalServerError)
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

		data := struct {
			Title      string
			User       db.User
			LoggedIn   bool
			Karma      int
			UserPosts  []db.Post
			LikedPosts []db.Post
		}{
			Title:      "Profile Page",
			LoggedIn:   true,
			User:       user,
			Karma:      karma,
			UserPosts:  userPosts,
			LikedPosts: likedPosts,
		}

		utils.RenderTemplate(w, "profile/profile.html", data)
	}
}
