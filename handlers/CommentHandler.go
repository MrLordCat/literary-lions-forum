package handlers

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

// Функция для добавления комментария в базу данных
func AddComment(db *sql.DB, postID, authorID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (post_id, author_id, content) VALUES (?, ?, ?)", postID, authorID, content)
	fmt.Println(content)
	return err

}

func AddCommentHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := GetUserIDFromSession(r)
		if err != nil || userID == 0 {
			http.Error(w, "You must be logged in to post comments", http.StatusForbidden)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
			return
		}

		err = AddComment(dbConn, postID, userID, content)
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			return
		}

		// Получаем ID автора поста
		authorPosts, err := db.GetAllPosts(dbConn, int64(postID), 0)
		if err != nil {
			http.Error(w, "Failed to get post author", http.StatusInternalServerError)
			return
		}

		if len(authorPosts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		authorID := authorPosts[0].AuthorID

		// Создаем уведомление для автора поста
		contentNotification := "Your post has been commented on."
		err = db.AddNotification(dbConn, authorID, contentNotification)
		if err != nil {
			http.Error(w, "Failed to add notification", http.StatusInternalServerError)
			return
		}

		// Перенаправление обратно на страницу постов
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
