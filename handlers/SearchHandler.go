package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
)

func SearchHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		results, err := db.SearchAll(dbConn, query)
		if err != nil {
			http.Error(w, "Failed to fetch search results", http.StatusInternalServerError)
			return
		}

		loggedInUserID, _ := GetUserIDFromSession(r)
		isAdmin := false
		if loggedInUserID > 0 {
			user, err := db.GetUserByID(dbConn, loggedInUserID)
			if err == nil && user.IsAdmin {
				isAdmin = true
			}
		}

		options := map[string]bool{
			"notifications": true,
		}
		pageData, err := utils.GetPageData(dbConn, loggedInUserID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Results":             results,
			"Title":               "Search Results",
			"LoggedIn":            true,
			"IsAdmin":             isAdmin,
			"Notifications":       pageData.Notifications,
			"UnreadNotifications": pageData.UnreadNotifications,
		}

		// Подготовка данных для шаблона posts
		postsData := struct {
			IsProfile     bool
			IsAdmin       bool
			Notifications []db.Notification
			Posts         []db.Post
		}{
			IsProfile:     false,
			IsAdmin:       isAdmin,
			Posts:         results.Posts,
			Notifications: pageData.Notifications,
		}

		// Включение данных для шаблона posts
		data["PostsData"] = postsData

		utils.RenderTemplate(w, "searchResults.html", data)
	}
}
