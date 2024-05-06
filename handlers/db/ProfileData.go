package db

import (
	"database/sql"
	"log"
)

func GetUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, first_name, last_name FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
	return user, err
}
func GetUserPosts(db *sql.DB, userID int) ([]Post, error) {
	var posts []Post
	log.Printf("Fetching posts for user ID: %d", userID)
	rows, err := db.Query("SELECT id, title, content, created_at, category_id FROM posts WHERE author_id = ?", userID)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.CategoryID)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			return nil, err
		}
		posts = append(posts, p)
	}
	log.Printf("Retrieved %d posts", len(posts))
	return posts, nil
}

func GetLikedPosts(db *sql.DB, userID int) ([]Post, error) {
	var posts []Post
	query := `
SELECT p.id, p.title, p.content, p.created_at, p.category_id
FROM posts p
JOIN post_likes pl ON p.id = pl.post_id
WHERE pl.user_id = ? AND pl.like_type = 1
ORDER BY p.created_at DESC
`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.CategoryID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
