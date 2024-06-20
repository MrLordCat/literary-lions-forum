package handlers

import (
	"database/sql"
	"encoding/json"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

// Функция для добавления комментария в базу данных
func AddComment(db *sql.DB, postID, authorID int, content string) (int, error) {
	result, err := db.Exec("INSERT INTO comments (post_id, author_id, content) VALUES (?, ?, ?)", postID, authorID, content)
	if err != nil {
		return 0, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentID), nil
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

		commentID, err := AddComment(dbConn, postID, userID, content)
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			return
		}

		// Получаем новый комментарий
		newComment, err := db.GetCommentByID(dbConn, commentID)
		if err != nil {
			http.Error(w, "Failed to fetch new comment", http.StatusInternalServerError)
			return
		}

		// Создаем уведомление для автора поста
		authorPosts, err := db.GetAllPosts(dbConn, postID, 0, "likes")
		if err != nil {
			http.Error(w, "Failed to get post author", http.StatusInternalServerError)
			return
		}

		if len(authorPosts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		authorID := authorPosts[0].AuthorID
		contentNotification := "Your post has been commented on."
		err = db.AddNotification(dbConn, authorID, postID, contentNotification)
		if err != nil {
			http.Error(w, "Failed to add notification", http.StatusInternalServerError)
			return
		}

		// Возвращаем JSON-ответ с новым комментарием
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newComment)
	}
}
