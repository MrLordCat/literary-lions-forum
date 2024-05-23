package handlers

import (
	"database/sql"
	"fmt"
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
		fmt.Println(query)
		results, err := db.SearchAll(dbConn, query)
		if err != nil {
			http.Error(w, "Failed to fetch search results", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Results":  results,
			"Title":    "Search Results",
			"LoggedIn": true,
		}

		utils.RenderTemplate(w, "searchResults.html", data)
	}
}
