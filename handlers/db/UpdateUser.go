package db

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(db *sql.DB, userID int, username, firstName, lastName, password string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec("UPDATE users SET username = ?, first_name = ?, last_name = ? WHERE id = ?", username, firstName, lastName, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
