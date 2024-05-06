package db

import "database/sql"

type User struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    string
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
type Post struct {
	ID         int
	Title      string
	Content    string
	AuthorID   sql.NullString
	AuthorName string
	CreatedAt  string
	CategoryID sql.NullString
	Comments   []Comment
	Likes      int
	Dislikes   int
}
type Comment struct {
	AuthorName string
	Content    string
	CreatedAt  string
}
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
