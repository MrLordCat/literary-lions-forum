package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func CreateCategoryHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if r.FormValue("action") == "create" {
			categoryName := r.FormValue("categoryName")
			if categoryName == "" {
				http.Error(w, "Category name must not be empty", http.StatusBadRequest)
				return
			}

			
			if err := db.CreateCategory(dbConn, categoryName); err != nil {
				http.Error(w, "Failed to create category", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
		} else if r.FormValue("action") == "delete" {
			categoryName := r.FormValue("categoryNameToDelete")
			if categoryName == "" {
				http.Error(w, "Category name must not be empty", http.StatusBadRequest)
				return
			}

			
			if err := db.DeleteCategory(dbConn, categoryName); err != nil {
				http.Error(w, "Failed to delete category", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
