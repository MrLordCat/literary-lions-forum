package handlers

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"strconv"
)

func PostsHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.GetAllPosts(dbConn, 0, 0)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("web/template/posts.html"))
		err = tmpl.Execute(w, map[string]interface{}{"Posts": posts})
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}
func PostViewHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.URL.Query().Get("postID") // Получение ID поста из URL параметра
		if postIDStr == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		posts, err := db.GetAllPosts(dbConn, postID, 0)
		if err != nil {
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
			return
		}
		if len(posts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		post := posts[0]

		comments, err := db.GetCommentsForPost(dbConn, int(postID))
		if err != nil {
			http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Post":     post,
			"Comments": comments,
			"LoggedIn": true,
		}

		utils.RenderTemplate(w, "post/postView.html", data)
	}
}
