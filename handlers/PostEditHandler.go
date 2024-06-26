package handlers

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"strconv"
)

func EditPostHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.FormValue("postID")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		posts, err := db.GetAllPosts(dbConn, postID, 0, "likes")
		if err != nil {
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
			return
		}
		if len(posts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		post := posts[0]

		
		canEdit := false

		loggedInUserID, err := GetUserIDFromSession(r)
		if err == nil {
			user, err := db.GetUserByID(dbConn, loggedInUserID)
			if err == nil && (user.ID == post.AuthorID || user.IsAdmin) {
				canEdit = true
			}
		}

		options := map[string]bool{
			"notifications": true,
		}

		pageData, err := utils.GetPageData(dbConn, loggedInUserID, options)
		if err != nil {
			http.Error(w, "Failed to fetch page data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageData.Post = post
		pageData.CanEdit = canEdit
		pageData.Title = "Edit Post"

		utils.RenderTemplate(w, "post/editPost.html", pageData)
	}
}
func UpdatePostHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		postIDStr := r.FormValue("postID")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			fmt.Print(postID)
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromSession(r)
		if err != nil || userID == 0 {
			http.Error(w, "You must be logged in", http.StatusForbidden)
			return
		}

		user, err := db.GetUserByID(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		posts, err := db.GetAllPosts(dbConn, postID, 0, "likes")
		if err != nil {
			http.Error(w, "Failed to get post", http.StatusInternalServerError)
			return
		}

		if len(posts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		post := posts[0]

		if user.ID != post.AuthorID && !user.IsAdmin {
			http.Error(w, "You do not have permission to edit this post", http.StatusForbidden)
			return
		}

		if r.FormValue("action") == "delete" {
			
			if err := db.UpdateOrDeletePost(dbConn, postID, "", "", "", "", "", true); err != nil {
				http.Error(w, "Failed to delete post", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		
		title := r.FormValue("title")
		content := r.FormValue("content")

		
		image1Path, err := SaveUploadedFile(r, "image1")
		if err != nil {
			http.Error(w, "Failed to save image 1", http.StatusInternalServerError)
			return
		}
		image2Path, err := SaveUploadedFile(r, "image2")
		if err != nil {
			http.Error(w, "Failed to save image 2", http.StatusInternalServerError)
			return
		}
		image3Path, err := SaveUploadedFile(r, "image3")
		if err != nil {
			http.Error(w, "Failed to save image 3", http.StatusInternalServerError)
			return
		}

		
		if r.FormValue("delete_image1") == "on" {
			image1Path = "DELETE"
		}
		if r.FormValue("delete_image2") == "on" {
			image2Path = "DELETE"
		}
		if r.FormValue("delete_image3") == "on" {
			image3Path = "DELETE"
		}

		if err := db.UpdateOrDeletePost(dbConn, postID, title, content, image1Path, image2Path, image3Path, false); err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}
		fmt.Println(content)
		http.Redirect(w, r, "/postView?postID="+strconv.Itoa(int(postID)), http.StatusFound)
	}
}
