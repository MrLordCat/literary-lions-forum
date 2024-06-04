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

		// Обработка контента каждого поста
		for i, post := range posts {
			posts[i].Content = string(utils.RenderPostContent(post.Content))
		}

		tmpl := template.Must(template.New("posts.html").Funcs(template.FuncMap{
			"renderPostContent": utils.RenderPostContent,
		}).ParseFiles("web/templates/post/posts.html"))

		err = tmpl.Execute(w, map[string]interface{}{"Posts": posts})
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func PostViewHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.URL.Query().Get("postID")
		if postIDStr == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to view this page", http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(dbConn, userID)
		if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
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
			"UserID":   userID,
			"IsAdmin":  user.IsAdmin,
		}

		utils.RenderTemplate(w, "post/postView.html", data)
	}
}
