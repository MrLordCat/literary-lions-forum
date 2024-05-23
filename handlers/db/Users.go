package db

import (
	"database/sql"
	"errors"
	"log"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, username, email, password string) error {
	// Проверка формата email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, string(hashedPassword))
	return err
}

func GetUserByUsernameOrEmail(db *sql.DB, login string) (User, error) {
	var user User
	// Обновлённый SQL запрос для поиска по username или email
	query := "SELECT id, username, email, password_hash FROM users WHERE username = ? OR email = ?"
	err := db.QueryRow(query, login, login).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
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
func CheckCurrentPassword(db *sql.DB, userID int, currentPassword string) bool {
	var hashedPassword string
	err := db.QueryRow(`SELECT password_hash FROM users WHERE id = ?`, userID).Scan(&hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	return err == nil
}
