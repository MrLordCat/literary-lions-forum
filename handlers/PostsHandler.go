package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"strconv"
)

/*
	 func PostsHandler(dbConn *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			

			userID, err := GetUserIDFromSession(r)
			if err != nil {
				http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
				return
			}
			
			options := map[string]bool{
				"posts":         true,
				"notifications": true,
				"IsAdmin":       true,
			}

			
			pageData, err := utils.GetPageData(dbConn, userID, options)
			if err != nil {
				http.Error(w, "Failed to fetch page data", http.StatusInternalServerError)
				return
			}

			
			err = utils.RenderTemplate(w, "posts.html", pageData)
			if err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
		}
	}
*/
func PostViewHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		postIDStr := r.URL.Query().Get("postID")
		if postIDStr == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		userID, _ := GetUserIDFromSession(r)

		options := map[string]bool{
			"singlePost": true,
		}

		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		posts, err := db.GetAllPosts(dbConn, postID, 0, "likes")
		if err != nil {
			http.Error(w, "Failed to fetch post: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(posts) > 0 {
			data.Posts = posts
		}

		utils.RenderTemplate(w, "post/postView.html", data)
	}
}
