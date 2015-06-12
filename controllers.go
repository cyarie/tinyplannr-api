package main

import (
	"database/sql"
	"log"
)

func getUserById(db *sql.DB, id int64) (*User, error) {
	var retval User

	query_str, err := db.Prepare(`SELECT user_id, email, first_name, last_name, zip_code, is_active, create_dt, update_dt
	                              FROM tinyplannr_api.user_api
	                              WHERE user_id = $1`)
	if err != nil {
		panic(err)
	}

	err = query_str.QueryRow(id).Scan(&retval.ID, &retval.Email, &retval.FirstName, &retval.LastName,
		&retval.ZipCode, &retval.IsActive, &retval.CreateDt, &retval.UpdateDt)
	if err != nil {
		return &retval, err
	}
	retval.ID = id

	return &retval, err
}

func createUserDb(db *sql.DB, u User) (*User, error) {

	query_str, err := db.Prepare(`INSERT INTO tinyplannr_api.user_api
	                                  VALUES (DEFAULT, $1, $2, $3, $4, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	                                  RETURNING user_id`)
	if err != nil {
		log.Fatal(err)
	}

	err = query_str.QueryRow(u.Email, u.FirstName, u.LastName, u.ZipCode).Scan(&u.ID)
	if err != nil {
		panic(err)
	}

	lastId := u.ID

	retval, err := getUserById(db, lastId)
	if err != nil {
		log.Fatal(err)
	}

	return retval, err
}
