package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreatePost(db *sql.DB, title, content string, authorID, categoryID int, image1Path, image2Path, image3Path string) error {
	_, err := db.Exec(`
		INSERT INTO posts (title, content, author_id, category_id, image1_path, image2_path, image3_path, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`, title, content, authorID, categoryID, image1Path, image2Path, image3Path)
	return err
}

func GetAllPosts(db *sql.DB, postID int, userID int64, sortBy string) ([]Post, error) {
	var query strings.Builder
	query.WriteString(`
	SELECT p.id, p.title, p.content, u.username, u.id as author_id, p.created_at, p.is_deleted,
	COALESCE(SUM(CASE WHEN pl.like_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
	COALESCE(SUM(CASE WHEN pl.like_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
	p.image1_path, p.image2_path, p.image3_path,
	c.name as category_name
	FROM posts p
	LEFT JOIN users u ON p.author_id = u.id
	LEFT JOIN post_likes pl ON p.id = pl.post_id
	LEFT JOIN categories c ON p.category_id = c.id
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

	query.WriteString(" GROUP BY p.id, p.title, p.content, u.username, u.id, p.created_at, p.image1_path, p.image2_path, p.image3_path, c.name")

	switch sortBy {
	case "likes":
		query.WriteString(" ORDER BY likes DESC")
	case "date":
		fallthrough
	default:
		query.WriteString(" ORDER BY p.created_at DESC")
	}

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
		var image1Path, image2Path, image3Path sql.NullString
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &authorName, &authorID, &p.CreatedAt, &p.IsDeleted, &p.Likes, &p.Dislikes, &image1Path, &image2Path, &image3Path, &p.CategoryName)
		if err != nil {
			fmt.Println("Error scanning post:", err)
			continue
		}
		p.AuthorName = authorName.String
		if authorID.Valid {
			p.AuthorID = int(authorID.Int64)
		}
		if image1Path.Valid {
			p.Image1Path = image1Path.String
		}
		if image2Path.Valid {
			p.Image2Path = image2Path.String
		}
		if image3Path.Valid {
			p.Image3Path = image3Path.String
		}
		p.Comments, _ = GetCommentsForPost(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostsByCategory(db *sql.DB, categoryID int, sort string) ([]Post, error) {
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
	`

	switch sort {
	case "likes":
		query += " ORDER BY likes DESC"
	case "date":
		fallthrough
	default:
		query += " ORDER BY p.created_at DESC"
	}

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

func UpdateOrDeletePost(db *sql.DB, postID int, title, content, image1Path, image2Path, image3Path string, delete bool) error {
	if delete {
		
		query := `UPDATE posts SET title = '', content = '', is_deleted = 1 WHERE id = ?`
		_, err := db.Exec(query, postID)
		if err != nil {
			return fmt.Errorf("error marking post as deleted: %v", err)
		}
	} else {
		
		fmt.Println(image1Path)
		query := `UPDATE posts SET title = ?, content = ?, 
            image1_path = CASE WHEN ? = 'DELETE' THEN NULL WHEN ? != '' THEN ? ELSE image1_path END, 
            image2_path = CASE WHEN ? = 'DELETE' THEN NULL WHEN ? != '' THEN ? ELSE image2_path END, 
            image3_path = CASE WHEN ? = 'DELETE' THEN NULL WHEN ? != '' THEN ? ELSE image3_path END 
            WHERE id = ?`
		_, err := db.Exec(query, title, content, image1Path, image1Path, image1Path, image2Path, image2Path, image2Path, image3Path, image3Path, image3Path, postID)
		if err != nil {
			return fmt.Errorf("error updating post: %v", err)
		}
	}
	return nil
}
