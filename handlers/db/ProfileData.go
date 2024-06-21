package db

import (
	"database/sql"
)

func GetUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, first_name, last_name, karma, is_admin FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Karma, &user.IsAdmin)
	return user, err
}

func GetLikedPosts(db *sql.DB, userID int) ([]Post, error) {
	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.is_deleted
              FROM post_likes pl
              JOIN posts p ON pl.post_id = p.id
              WHERE pl.user_id = ? AND pl.like_type = 1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.IsDeleted)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
