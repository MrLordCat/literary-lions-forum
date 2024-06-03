package db

import (
	"database/sql"
	"fmt"
)

func GetCommentByID(db *sql.DB, commentID int) (*Comment, error) {
	query := `
    SELECT c.id, c.content, c.author_id, cu.username, c.created_at, COALESCE(SUM(cl.like_type), 0) AS likes
    FROM comments c
    JOIN users cu ON c.author_id = cu.id
    LEFT JOIN comment_likes cl ON c.id = cl.comment_id
    WHERE c.id = ?
    GROUP BY c.id, c.content, c.author_id, cu.username, c.created_at
    `
	row := db.QueryRow(query, commentID)

	var c Comment
	err := row.Scan(&c.ID, &c.Content, &c.AuthorID, &c.AuthorName, &c.CreatedAt, &c.Likes)
	if err == sql.ErrNoRows {
		return nil, nil // No comment found
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func GetCommentsForPost(db *sql.DB, postID int) ([]Comment, error) {
	query := `
    SELECT c.id, c.content, cu.username, c.created_at, COALESCE(SUM(cl.like_type), 0) AS likes
    FROM comments c
    JOIN users cu ON c.author_id = cu.id
    LEFT JOIN comment_likes cl ON c.id = cl.comment_id
    WHERE c.post_id = ?
    GROUP BY c.id
    ORDER BY c.created_at DESC
    `
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.AuthorName, &c.CreatedAt, &c.Likes); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func DeleteOrUpdateComment(db *sql.DB, commentID int, newContent string, delete bool) error {
	if delete {
		// Если delete == true, удаляем комментарий
		query := `DELETE FROM comments WHERE id = ?`
		_, err := db.Exec(query, commentID)
		if err != nil {
			return fmt.Errorf("error deleting comment: %v", err)
		}
	} else {
		// Если delete == false, обновляем комментарий
		query := `UPDATE comments SET content = ? WHERE id = ?`
		_, err := db.Exec(query, newContent, commentID)
		if err != nil {
			return fmt.Errorf("error updating comment: %v", err)
		}
	}
	return nil
}
