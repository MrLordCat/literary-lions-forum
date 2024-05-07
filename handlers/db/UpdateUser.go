package db

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(db *sql.DB, userID int, username, firstName, lastName, password string) error {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Если пользователь указал новый пароль, обновите его.
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Предполагается, что bcrypt уже импортирован
		if err != nil {
			tx.Rollback() // Откат транзакции в случае ошибки
			return err
		}
		_, err = tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	fmt.Println("Sent username:", username, firstName, lastName)
	// Обновление остальной информации пользователя
	_, err = tx.Exec("UPDATE users SET username = ?, first_name = ?, last_name = ? WHERE id = ?", username, firstName, lastName, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Подтверждение транзакции
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
func СheckCurrentPassword(db *sql.DB, userID int, currentPassword string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&hashedPassword)
	if err != nil {
		log.Println("Error fetching user password:", err)
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword)) == nil
}
