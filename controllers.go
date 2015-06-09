package main

import (
	"database/sql"
)

func getUserById(db *sql.DB, id int) (*User, error) {
	const query = `SELECT id, email FROM users WHERE id = $1 `
	var retval User
	err := db.QueryRow(query, id).Scan(&retval.ID, &retval.Email)
	retval.ID = id
	return &retval, err
}
