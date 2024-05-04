package db

import (
	"database/sql"
	"fmt"
)

func CreatePost(db *sql.DB, title, content string, authorID int) error {
	stmt, err := db.Prepare("INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, authorID)
	return err
}

type Post struct {
	ID         int
	Title      string
	Content    string
	AuthorID   sql.NullString
	AuthorName string
	CreatedAt  string
	Comments   []Comment
	Likes      int // количество лайков
	Dislikes   int
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
	SELECT p.id, p.title, p.content, u.username, p.created_at,
	COALESCE(SUM(CASE WHEN pl.like_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
	COALESCE(SUM(CASE WHEN pl.like_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
FROM posts p
LEFT JOIN users u ON p.author_id = u.id
LEFT JOIN post_likes pl ON p.id = pl.post_id
GROUP BY p.id, p.title, p.content, u.username, p.created_at
ORDER BY p.created_at DESC
    `)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var authorName sql.NullString
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &authorName, &p.CreatedAt, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}
		p.AuthorName = authorName.String
		p.Comments, _ = GetCommentsForPost(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil
}
