package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

func EditCommentHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
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

		comment, err := db.GetCommentByID(dbConn, commentID) // Используем функцию для получения конкретного комментария
		if err != nil {
			http.Error(w, "Failed to get comment", http.StatusInternalServerError)
			return
		}

		if comment == nil {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}

		if user.ID != comment.AuthorID && !user.IsAdmin {
			http.Error(w, "You do not have permission to edit this comment", http.StatusForbidden)
			return
		}

		// Проверка, какое действие нужно выполнить
		action := r.FormValue("action") // Действие может быть 'delete' или 'update'
		switch action {
		case "delete":
			if err := db.DeleteOrUpdateComment(dbConn, commentID, "", true); err != nil {
				http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
				return
			}
		case "update":
			newContent := r.FormValue("content")
			if err := db.DeleteOrUpdateComment(dbConn, commentID, newContent, false); err != nil {
				http.Error(w, "Failed to update comment", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
