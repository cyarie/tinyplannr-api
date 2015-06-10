package main

import (
	"database/sql"
	"fmt"
)


func getUserById(db *sql.DB, id int) (*User, error) {
	var retval User

	query_str, err := db.Prepare("SELECT id, email, created_dt FROM users WHERE id = $1")
	if err != nil {
		panic(err)
	}

	fmt.Println("Querying the database...")
	err = query_str.QueryRow(id).Scan(&retval.ID, &retval.Email, &retval.CreatedDt)
	if err != nil {
		return &retval, err
	}
	fmt.Println("Query finished...")
	retval.ID = id

	return &retval, err
}
