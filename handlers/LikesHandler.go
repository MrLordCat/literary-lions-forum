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
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
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

		// Получаем автора поста для обновления его кармы
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

		// Обновляем карму автора поста
		err = db.UpdateUserKarma(dbConn, post.AuthorID)
		if err != nil {
			http.Error(w, "Failed to update user karma", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func LikeCommentHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := GetUserIDFromSession(r)
		if err != nil || userID == 0 {
			http.Error(w, "You must be logged in to like or dislike comments", http.StatusForbidden)
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

		likeType, err := strconv.Atoi(r.FormValue("like_type")) // 1 for like, -1 for dislike
		if err != nil || (likeType != 1 && likeType != -1) {
			http.Error(w, "Invalid like type", http.StatusBadRequest)
			return
		}

		if err := db.LikeComment(dbConn, commentID, userID, likeType); err != nil {
			http.Error(w, "Failed to like comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Получаем автора комментария для обновления его кармы
		comment, err := db.GetCommentByID(dbConn, commentID)
		if err != nil {
			http.Error(w, "Failed to fetch comment", http.StatusInternalServerError)
			return
		}

		// Обновляем карму автора комментария
		err = db.UpdateUserKarma(dbConn, comment.AuthorID)
		if err != nil {
			http.Error(w, "Failed to update user karma", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
