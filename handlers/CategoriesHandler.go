package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func CreateCategoryHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryName := r.FormValue("categoryName")
	if categoryName == "" {
		http.Error(w, "Category name must not be empty", http.StatusBadRequest)
		return
	}

	// Добавление категории в базу данных
	if err := db.CreateCategory(dbConn, categoryName); err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusFound)
}
