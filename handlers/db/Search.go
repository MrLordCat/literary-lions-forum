package db

import (
	"database/sql"
	"fmt"
	"log"
)

func SearchAll(db *sql.DB, query string) (SearchResult, error) {
	var result SearchResult

	// Поиск постов по названию или содержанию
	posts, err := SearchPosts(db, query)
	if err != nil {
		return result, err
	}
	result.Posts = posts

	// Поиск пользователей по имени пользователя или email
	users, err := SearchUsers(db, query)
	if err != nil {
		return result, err
	}
	result.Users = users

	return result, nil
}
func SearchPosts(db *sql.DB, searchText string) ([]Post, error) {
	// Используйте searchText для модификации SQL запроса
	posts := []Post{}
	query := `SELECT p.id, p.title, p.content, u.username AS author_name, p.author_id, p.created_at, p.category_id,
	COALESCE(SUM(CASE WHEN pl.like_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
	COALESCE(SUM(CASE WHEN pl.like_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
	FROM posts p
	LEFT JOIN users u ON p.author_id = u.id
	LEFT JOIN post_likes pl ON p.id = pl.post_id
	WHERE (p.title LIKE ? OR p.content LIKE ?) AND p.is_deleted = 0
	GROUP BY p.id, p.title, p.content, u.username, p.author_id, p.created_at
	ORDER BY p.created_at DESC
	`
	// Используйте '%' + searchText + '%' для LIKE условий
	rows, err := db.Query(query, "%"+searchText+"%", "%"+searchText+"%")
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
	fmt.Println("Result:", posts)
	return posts, nil
}

func SearchUsers(db *sql.DB, searchText string) ([]User, error) {
	// Используйте searchText для модификации SQL запроса
	query := `SELECT id, username, email, created_at FROM users
	WHERE LOWER(username) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)
	
	`
	rows, err := db.Query(query, "%"+searchText+"%", "%"+searchText+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("no rows in result set")
		}
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			log.Println("Failed to scan row: ", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}
