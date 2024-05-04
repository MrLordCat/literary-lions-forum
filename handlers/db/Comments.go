package db

import "database/sql"

type Comment struct {
	AuthorName string
	Content    string
	CreatedAt  string
}

func GetCommentsForPost(db *sql.DB, postID int) ([]Comment, error) {
	query := `
    SELECT c.content, cu.username, c.created_at
    FROM comments c
    JOIN users cu ON c.author_id = cu.id
    WHERE c.post_id = ?
    ORDER BY c.created_at
    `
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.Content, &c.AuthorName, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
