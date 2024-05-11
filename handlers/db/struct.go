package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Username     string
	Email        string
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
	IsDeleted  bool
}
type Comment struct {
	ID         int
	AuthorName string
	Content    string
	CreatedAt  string
	Likes      int // Добавляем поле для хранения количества лайков
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
