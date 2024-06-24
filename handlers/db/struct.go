package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Karma        sql.NullInt64
	Username     string
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
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
	CategoryName string 
	Comments   []Comment
	Likes      int
	Dislikes   int
	IsDeleted  bool
	Image1Path string
	Image2Path string
	Image3Path string
}
type Comment struct {
	ID         int
	AuthorName string
	AuthorID   int
	Content    string
	CreatedAt  time.Time
	Likes      int 
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type SearchResult struct {
	Posts []Post
	Users []User
}
type Notification struct {
	ID        int
	UserID    int
	Content   string
	IsRead    int
	CreatedAt string
	PostID    int
}
