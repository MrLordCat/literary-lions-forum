package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"sort"
)

func UsersHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
			return
		}

		var users []db.User
		users, err = db.GetAllUsers(dbConn)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		// По умолчанию сортируем по карме
		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" || sortBy == "karma" {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Karma > users[j].Karma
			})
		}

		options := map[string]bool{
			"notifications": true,
		}

		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Title = "Users"
		data.Users = users

		utils.RenderTemplate(w, "users.html", data)
	}
}
