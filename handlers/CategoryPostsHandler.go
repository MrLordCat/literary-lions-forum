package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"literary-lions-forum/handlers/db"
	"log"
	"net/http"
	"strconv"
)

func CategoryPostsHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIDStr := r.URL.Query().Get("category_id")
		if categoryIDStr == "" {
			http.Error(w, "Category ID is required", http.StatusBadRequest)
			return
		}
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var categoryName string
		err = db.QueryRow(dbConn, "SELECT name FROM categories WHERE id = ?", []interface{}{categoryID}, &categoryName)
		if err != nil {
			http.Error(w, "Failed to fetch category name: "+err.Error(), http.StatusInternalServerError)
			return
		}

		posts, err := db.GetPostsByCategory(dbConn, categoryID)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		fmt.Println(categoryName)
		tmpl := template.Must(template.ParseFiles("web/template/sortedPosts.html"))
		err = tmpl.Execute(w, map[string]interface{}{
			"Posts":        posts,
			"CategoryID":   categoryID,
			"CategoryName": categoryName, // Передаем имя категории в шаблон
		})
		log.Printf("Rendering posts for category: %s with ID %d", categoryName, categoryID)

		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}
