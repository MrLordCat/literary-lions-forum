package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)


func CreateUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func GetUserByUsername(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			log.Println("Failed to scan row: ", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}
