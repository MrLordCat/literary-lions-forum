package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func RegisterHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Устанавливаем isAdmin в true, если username равен 'admin'
		isAdmin := username == "admin"

		if err := db.CreateUser(dbConn, username, email, password, isAdmin); err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
