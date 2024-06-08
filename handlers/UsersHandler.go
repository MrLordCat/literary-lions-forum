package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"log"
	"net/http"
)

func UsersHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.GetAllUsers(dbConn)
		if err != nil {
			log.Printf("Error fetching users: %v", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Users []db.User
		}{
			Title: "Users",
			Users: users,
		}

		err = utils.RenderTemplate(w, "usersList.html", data)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
