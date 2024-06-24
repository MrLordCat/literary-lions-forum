package db

import (
	"database/sql"
)

func CalculateUserKarma(db *sql.DB, userID int) (int, error) {
	var postKarma int
	err := db.QueryRow(`
        SELECT COALESCE(SUM(like_type), 0) FROM post_likes WHERE user_id = ?
    `, userID).Scan(&postKarma)
	if err != nil {
		return 0, err
	}
	var commentKarma int
	err = db.QueryRow(`
        SELECT COALESCE(SUM(like_type), 0) FROM comment_likes WHERE user_id = ?
    `, userID).Scan(&commentKarma)
	if err != nil {
		return 0, err
	}
	commentKarma /= 2
	totalKarma := postKarma + commentKarma
	if totalKarma < 0 {
		totalKarma = 0
	}
	return totalKarma, nil
}
func UpdateUserKarma(db *sql.DB, userID int) error {
	var postLikes int
	postLikesQuery := `SELECT COALESCE(SUM(like_type), 0) FROM post_likes WHERE post_id IN (SELECT id FROM posts WHERE author_id = ?)`
	err := db.QueryRow(postLikesQuery, userID).Scan(&postLikes)
	if err != nil {
		return err
	}
	var commentLikes int
	commentLikesQuery := `SELECT COALESCE(SUM(like_type), 0) FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE author_id = ?)`
	err = db.QueryRow(commentLikesQuery, userID).Scan(&commentLikes)
	if err != nil {
		return err
	}
	karma := postLikes + (commentLikes / 2)
	updateQuery := `UPDATE users SET karma = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, karma, userID)
	if err != nil {
		return err
	}
	return nil
}
