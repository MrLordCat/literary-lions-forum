package db

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./handlers/db/literary_lions.db")
	if err != nil {
		log.Fatal(err)
	}

	CreateTables(db)
	return db
}

func CreateTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        author_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        category_id INTEGER,
        FOREIGN KEY (author_id) REFERENCES users(id),
        FOREIGN KEY (category_id) REFERENCES categories(id)
    );
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT DEFAULT '',
		last_name TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		author_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (author_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS comment_likes (
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		like_type INTEGER NOT NULL,  -- 1 для лайка, -1 для дизлайка
		PRIMARY KEY (comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS post_likes (
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		like_type INTEGER NOT NULL,  -- 1 для лайка, -1 для дизлайка
		PRIMARY KEY (post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT ''
	);

	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to execute query to create posts table: %s", err)

	}
	log.Println("Table posts created successfully or already exists.")

}

func SessionHandler(user User, w http.ResponseWriter) {
	// Создание сессионной куки
	sessionToken := "some_generated_session_token" // Реальная токен нужно генерировать безопасно
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Minute),
	})
}
