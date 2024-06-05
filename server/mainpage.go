package server

import (
	"database/sql"

	"literary-lions-forum/handlers"
	"literary-lions-forum/utils"
	"net/http"
)

func MainPageHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := handlers.GetUserIDFromSession(r)
		loggedIn := err == nil && userID != 0

		options := map[string]bool{
			"posts":         true,
			"notifications": loggedIn,
			"categories":    true,
			"topUsers":      true,
		}

		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Title = "Home"
		data.LoggedIn = loggedIn

		utils.RenderTemplate(w, "home/home.html", data)
	}
}
