package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
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

		var categoryID int
		var err error

		
		query := "SELECT id FROM categories WHERE name = ?"
		err = db.QueryRow(dbConn, query, []interface{}{categoryIDStr}, &categoryID)
		if err != nil {
			
			categoryID, err = strconv.Atoi(categoryIDStr)
			if err != nil {
				http.Error(w, "Invalid category ID or name", http.StatusBadRequest)
				return
			}
		}

		sort := r.URL.Query().Get("sort")
		if sort == "" {
			sort = "likes"
		}

		options := map[string]bool{
			"notifications": true,
			"categories":    true,
		}

		pageData, err := utils.GetPageData(dbConn, 0, options)
		if err != nil {
			http.Error(w, "Failed to fetch page data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var categoryName string
		err = db.QueryRow(dbConn, "SELECT name FROM categories WHERE id = ?", []interface{}{categoryID}, &categoryName)
		if err != nil {
			http.Error(w, "Failed to fetch category name: "+err.Error(), http.StatusInternalServerError)
			return
		}

		posts, err := db.GetPostsByCategory(dbConn, categoryID, sort)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		pageData.CategoryName = categoryName
		pageData.Posts = posts
		pageData.Title = "Sorted Posts by Category"
		pageData.Sort = sort
		pageData.CategoryID = categoryID 

		err = utils.RenderTemplate(w, "post/sortedPosts.html", pageData)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
