// Storing all of our database controller functions in this file.

package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost	int = 13
)

func getUserById(db *sql.DB, id int64) (*UserDisplay, error) {
	var retval UserDisplay

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

func createUserAuth(db *sql.DB, u UserCreate) {

	query_str, err := db.Prepare(`INSERT INTO tinyplannr_api.user_auth
	                                  VALUES (DEFAULT, $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	                                  RETURNING user_id`)

	if err != nil {
		log.Fatal(err)
	}

	password := []byte(u.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 13)
	if err != nil {
		log.Fatal(err)
	}

	err = query_str.QueryRow(u.ID, u.Email, string(hashedPassword)).Scan(&u.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created a User with password hash of: %s", string(hashedPassword))

}

func createUserDb(db *sql.DB, u UserCreate) (*UserDisplay, error) {

	query_str, err := db.Prepare(`INSERT INTO tinyplannr_api.user_api
	                                  VALUES (DEFAULT, $1, $2, $3, $4, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	                                  RETURNING user_id`)
	if err != nil {
		log.Fatal(err)
	}

	// Let's run the first query, which will create the publicly viewable user data
	err = query_str.QueryRow(u.Email, u.FirstName, u.LastName, u.ZipCode).Scan(&u.ID)
	if err != nil {
		panic(err)
	}

	// Now, let's create the UserAuth entry, which stores the password hash
	createUserAuth(db, u)

	lastId := u.ID

	retval, err := getUserById(db, lastId)
	if err != nil {
		log.Fatal(err)
	}

	return retval, err
}
