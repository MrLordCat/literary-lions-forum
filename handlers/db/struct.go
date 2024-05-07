package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    string
	FirstName    sql.NullString
	LastName     sql.NullString
}
type Post struct {
	ID         int
	Title      string
	Content    string
	AuthorID   int
	AuthorName string
	CreatedAt  time.Time
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
