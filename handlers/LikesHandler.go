package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

func LikeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	err = AddOrUpdateLike(db, postID, userID, likeType)
	if err != nil {
		http.Error(w, "Failed to update like/dislike", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusFound)
}

func AddOrUpdateLike(db *sql.DB, postID, userID, likeType int) error {
	// Check if the user has already liked or disliked the post
	var existingType int
	err := db.QueryRow("SELECT like_type FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		// No existing like/dislike, insert new
		_, err = db.Exec("INSERT INTO post_likes (post_id, user_id, like_type) VALUES (?, ?, ?)", postID, userID, likeType)
	} else if existingType != likeType {
		// Existing like/dislike differs, update it
		_, err = db.Exec("UPDATE post_likes SET like_type = ? WHERE post_id = ? AND user_id = ?", likeType, postID, userID)
	}
	return err
}
