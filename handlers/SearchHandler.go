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
			"IsAdmin":       isAdmin,
		}
		pageData, err := utils.GetPageData(dbConn, loggedInUserID, options)
		if err != nil {
			http.Error(w, "Failed to fetch page data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageData.Title = "Search Results"
		pageData.Posts = results.Posts

		utils.RenderTemplate(w, "searchResults.html", pageData)
	}
}
