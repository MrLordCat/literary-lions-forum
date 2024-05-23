package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

func LikeHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := GetUserIDFromSession(r)
		if err != nil || userID == 0 {
			http.Error(w, "You must be logged in to like or dislike posts", http.StatusForbidden)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		likeType, err := strconv.Atoi(r.FormValue("like_type")) // 1 for like, -1 for dislike
		if err != nil || (likeType != 1 && likeType != -1) {
			http.Error(w, "Invalid like type", http.StatusBadRequest)
			return
		}

		err = db.AddOrUpdateLike(dbConn, postID, userID, likeType)
		if err != nil {
			http.Error(w, "Failed to update like/dislike", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusFound)
	}
}
func LikeCommentHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		authorID, err := strconv.Atoi(r.FormValue("author_id")) // Используйте author_id здесь
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}

		likeType, err := strconv.Atoi(r.FormValue("like_type"))
		if err != nil {
			http.Error(w, "Invalid like type", http.StatusBadRequest)
			return
		}

		if err := db.LikeComment(dbConn, commentID, authorID, likeType); err != nil {
			http.Error(w, "Failed to like comment: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
