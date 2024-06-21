package server

import (
	"database/sql"

	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
)

func MainPageHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := handlers.GetUserIDFromSession(r)
		loggedIn := err == nil && userID != 0

		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			sortBy = "p.created_at DESC" // Сортировка по умолчанию
		}

		options := map[string]bool{
			"posts":         true,
			"notifications": loggedIn,
			"categories":    true,
			"topUsers":      true,
		}

		pageData, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageData.Posts, err = db.GetAllPosts(dbConn, 0, 0, sortBy)
		if err != nil {
			http.Error(w, "Failed to fetch posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageData.Title = "Home"
		pageData.LoggedIn = loggedIn
		pageData.Sort = sortBy
		utils.RenderTemplate(w, "home/home.html", pageData)
	}
}
