package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreatePost(db *sql.DB, title, content string, authorID, categoryID int) error {
	stmt, err := db.Prepare("INSERT INTO posts (title, content, author_id, category_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, authorID, categoryID)
	return err
}

func GetAllPosts(db *sql.DB, postID int64, userID int64) ([]Post, error) {
	var query strings.Builder
	query.WriteString(`
	SELECT p.id, p.title, p.content, u.username, u.id as author_id, p.created_at, p.is_deleted,
	COALESCE(SUM(CASE WHEN pl.like_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
	COALESCE(SUM(CASE WHEN pl.like_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
	FROM posts p
	LEFT JOIN users u ON p.author_id = u.id
	LEFT JOIN post_likes pl ON p.id = pl.post_id
	
    `)

	var args []interface{}
	whereClauses := []string{}
	if postID != 0 {
		whereClauses = append(whereClauses, "p.id = ?")
		args = append(args, postID)
	}
	if userID != 0 {
		whereClauses = append(whereClauses, "u.id = ?")
		args = append(args, userID)
	}

	if len(whereClauses) > 0 {
		query.WriteString(" WHERE " + strings.Join(whereClauses, " AND "))
	}

	query.WriteString(" GROUP BY p.id, p.title, p.content, u.username, u.id, p.created_at ORDER BY likes DESC")

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var authorName sql.NullString
		var authorID sql.NullInt64
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &authorName, &authorID, &p.CreatedAt, &p.IsDeleted, &p.Likes, &p.Dislikes)
		if err != nil {
			fmt.Println("Error scanning post:", err)
			continue
		}
		p.AuthorName = authorName.String
		if authorID.Valid {
			p.AuthorID = int(authorID.Int64)
		}
		p.Comments, _ = GetCommentsForPost(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil

}

func GetPostsByCategory(db *sql.DB, categoryID int64) ([]Post, error) {
	posts := []Post{}
	query := `
	SELECT p.id, p.title, p.content, u.username as author_name, p.author_id, p.created_at, p.category_id,
	COALESCE(SUM(CASE WHEN pl.like_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
	COALESCE(SUM(CASE WHEN pl.like_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
	FROM posts p
	LEFT JOIN users u ON p.author_id = u.id
	LEFT JOIN post_likes pl ON p.id = pl.post_id
	WHERE p.category_id = ? AND p.is_deleted = 0
	GROUP BY p.id, p.title, p.content, u.username, p.author_id, p.created_at
	ORDER BY p.created_at DESC
	`

	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		var authorID sql.NullInt64
		var authorName sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &authorName, &authorID, &p.CreatedAt, &p.CategoryID, &p.Likes, &p.Dislikes); err != nil {
			return nil, err
		}
		p.AuthorName = authorName.String
		if authorID.Valid {
			p.AuthorID = int(authorID.Int64)
		}
		p.Comments, _ = GetCommentsForPost(db, p.ID) 
		posts = append(posts, p)
	}

	return posts, nil
}

func UpdateOrDeletePost(db *sql.DB, postID int64, title, content string, delete bool) error {
	if delete {
		// Устанавливаем пост как удаленный и очищаем содержимое
		query := `UPDATE posts SET title = '', content = '', is_deleted = 1 WHERE id = ?`
		_, err := db.Exec(query, postID)
		if err != nil {
			return fmt.Errorf("error marking post as deleted: %v", err)
		}
	} else {
		// Обновляем пост с новым заголовком и содержимым
		query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
		_, err := db.Exec(query, title, content, postID)
		if err != nil {
			return fmt.Errorf("error updating post: %v", err)
		}
	}
	return nil
}
