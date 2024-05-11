package handlers

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
	"time"
)

func EditPostHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	posts, err := db.GetAllPosts(dbConn, postID, 0) // Предполагается, что первый параметр — это ID поста
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	if len(posts) == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	post := posts[0] // Получаем первый (единственный ожидаемый) пост
	if time.Since(post.CreatedAt) > time.Hour {
		http.Error(w, "You can only edit posts within one hour of creation", http.StatusForbidden)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/template/editPost.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"Post": post,
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func UpdatePostHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.ParseInt(r.FormValue("postID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if r.FormValue("action") == "delete" {
		// Обрабатываем удаление поста
		if err := db.UpdateOrDeletePost(dbConn, postID, "", "", true); err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Обрабатываем обновление поста
	title := r.FormValue("title")
	content := r.FormValue("content")
	if err := db.UpdateOrDeletePost(dbConn, postID, title, content, false); err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/postView?postID="+strconv.Itoa(int(postID)), http.StatusFound)
}
