package db

import (
	"database/sql"
	"fmt"
)

// QueryRow executes a query that is expected to return at most one row.
// query is the SQL query string, args are the arguments for any placeholders in the query,
// and dest are the pointers to variables where the scanned results should be stored.
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
