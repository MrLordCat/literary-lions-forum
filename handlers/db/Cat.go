package db

import (
	"database/sql"
)

func CreateCategory(dbConn *sql.DB, name string) error {
	_, err := dbConn.Exec("INSERT INTO categories (name) VALUES (?)", name)
	return err
}
func GetAllCategories(db *sql.DB) ([]Category, error) {
	categories := []Category{}
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
