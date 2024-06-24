package handlers

import (
	"database/sql" 
	"fmt"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"strconv"
)

func PostCreateFormHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			
			userID, err := GetUserIDFromSession(r)
			if err != nil {
				http.Error(w, "You need to be logged in to create a post", http.StatusUnauthorized)
				return
			}

			options := map[string]bool{
				"categories": true,
			}
			data, err := utils.GetPageData(dbConn, userID, options)
			if err != nil {
				http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
				return
			}

			utils.RenderTemplate(w, "/post/createPost.html", data)

		case "POST":
			
			userID, err := GetUserIDFromSession(r)
			if err != nil {
				http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
				return
			}
			if userID == 0 {
				http.Error(w, "User not logged in", http.StatusForbidden)
				return
			}

			
			title := r.FormValue("title")
			content := r.FormValue("content")
			categoryIDString := r.FormValue("category_id")
			categoryID, err := strconv.Atoi(categoryIDString)
			if err != nil {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			
			imagePaths := make([]string, 3)
			for i := 0; i < 3; i++ {
				formFile := fmt.Sprintf("image%d", i+1)
				imagePath, err := SaveUploadedFile(r, formFile)
				if err == nil {
					imagePaths[i] = imagePath
				}
			}

			
			if err := db.CreatePost(dbConn, title, content, userID, categoryID, imagePaths[0], imagePaths[1], imagePaths[2]); err != nil {
				http.Error(w, "Failed to save post: "+err.Error(), http.StatusInternalServerError)
				return
			}
			var postID int
			err = dbConn.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
			if err != nil {
				http.Error(w, "Failed to retrieve post ID: "+err.Error(), http.StatusInternalServerError)
				return
			}
			
			http.Redirect(w, r, "/postView?postID="+strconv.Itoa(postID), http.StatusFound)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
