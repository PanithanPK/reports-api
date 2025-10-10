package utils

import (
	"database/sql"
	"log"
)

// ExecWithLog executes a query and logs errors
func ExecWithLog(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Database exec error: %v | Query: %s", err, query)
	}
	return result, err
}

// QueryRowWithLog executes a single-row query and logs errors
func QueryRowWithLog(db *sql.DB, query string, args ...interface{}) *sql.Row {
	row := db.QueryRow(query, args...)
	log.Printf("Executing query: %s with args: %v", query, args)
	return row
}

// QueryWithLog executes a multi-row query and logs errors
func QueryWithLog(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Database query error: %v | Query: %s", err, query)
	}
	return rows, err
}

// ScanSingleValue scans a single value from a query
func ScanSingleValue(db *sql.DB, query string, dest interface{}, args ...interface{}) error {
	err := db.QueryRow(query, args...).Scan(dest)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Scan error: %v | Query: %s", err, query)
	}
	return err
}
