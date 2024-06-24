package db

import (
	"database/sql"
	"fmt"
)




func QueryRow(db *sql.DB, query string, args []interface{}, dest ...interface{}) error {
	if args == nil {
		args = []interface{}{}
	}
	row := db.QueryRow(query, args...)
	err := row.Scan(dest...)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no rows in result set")
		}
		return fmt.Errorf("query row error: %w", err)
	}
	return nil
}
