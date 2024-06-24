package db

import "database/sql"

func AddOrUpdateLike(db *sql.DB, postID, userID, likeType int) error {
	var existingType int
	err := db.QueryRow("SELECT like_type FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO post_likes (post_id, user_id, like_type) VALUES (?, ?, ?)", postID, userID, likeType)
	} else if existingType == likeType {
		_, err = db.Exec("DELETE FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
	} else {
		_, err = db.Exec("UPDATE post_likes SET like_type = ? WHERE post_id = ? AND user_id = ?", likeType, postID, userID)
	}
	return err
}

func LikeComment(db *sql.DB, commentID, userID, likeType int) error {
	var existingType int
	err := db.QueryRow("SELECT like_type FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&existingType)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO comment_likes (comment_id, user_id, like_type) VALUES (?, ?, ?)", commentID, userID, likeType)
		return err
	} else if err != nil {
		return err
	}
	if existingType == likeType {
		_, err = db.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		return err
	}
	_, err = db.Exec("UPDATE comment_likes SET like_type = ? WHERE comment_id = ? AND user_id = ?", likeType, commentID, userID)
	return err
}
