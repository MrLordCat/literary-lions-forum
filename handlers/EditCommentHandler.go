package handlers

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

func EditCommentHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID, err := strconv.ParseInt(r.FormValue("comment_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Проверка, какое действие нужно выполнить
	action := r.FormValue("action") // Действие может быть 'delete' или 'update'
	fmt.Println(action)
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
