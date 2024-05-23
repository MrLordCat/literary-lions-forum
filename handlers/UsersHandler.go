package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
	"text/template"
)

func UsersHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.GetAllUsers(dbConn)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("web/template/users.html"))
		tmpl.Execute(w, users)
	}
}
