package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"literary-lions-forum/handlers/db"
	"net/http"
	"strconv"
)

func handleLike(dbConn *sql.DB, r *http.Request, likeTypeKey, idKey string, isComment bool) (interface{}, error) {
	if r.Method != "POST" {
		return nil, fmt.Errorf("only POST method is allowed")
	}

	userID, err := GetUserIDFromSession(r)
	if err != nil || userID == 0 {
		return nil, fmt.Errorf("you must be logged in to like or dislike")
	}

	id, err := strconv.Atoi(r.FormValue(idKey))
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	likeType, err := strconv.Atoi(r.FormValue(likeTypeKey)) // 1 for like, -1 for dislike
	if err != nil || (likeType != 1 && likeType != -1) {
		return nil, fmt.Errorf("invalid like type")
	}

	if isComment {
		if err := db.LikeComment(dbConn, id, userID, likeType); err != nil {
			return nil, fmt.Errorf("failed to like comment: %v", err)
		}
		comment, err := db.GetCommentByID(dbConn, id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch comment: %v", err)
		}
		err = db.UpdateUserKarma(dbConn, comment.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user karma: %v", err)
		}
		return comment, nil
	} else {
		if err := db.AddOrUpdateLike(dbConn, id, userID, likeType); err != nil {
			return nil, fmt.Errorf("failed to update like/dislike: %v", err)
		}
		posts, err := db.GetAllPosts(dbConn, id, 0, "likes")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch post: %v", err)
		}
		if len(posts) == 0 {
			return nil, fmt.Errorf("post not found")
		}
		post := posts[0]
		err = db.UpdateUserKarma(dbConn, post.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user karma: %v", err)
		}
		return post, nil
	}
}

func LikeHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := handleLike(dbConn, r, "like_type", "post_id", false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post := result.(db.Post)

		// Возвращаем JSON ответ
		response := map[string]interface{}{
			"success":  true,
			"likes":    post.Likes,
			"dislikes": post.Dislikes,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func LikeCommentHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := handleLike(dbConn, r, "like_type", "comment_id", true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		comment := result.(*db.Comment)

		// Возвращаем JSON ответ
		response := map[string]interface{}{
			"success": true,
			"likes":   comment.Likes,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
