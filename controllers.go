// Storing all of our database controller functions in this file. Naming follows the pattern of:
// [function]Db or [handler]Db. These are mostly functions called by our HTTP handlers to interact with the database.
// I like keeping the SQL/DB trickery out of the handlers -- feels like a good separation of concerns.

package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost	int = 13
)

func getUserDb(db *sql.DB, id int64) (*UserDisplay, error) {
	var retval UserDisplay

	query_str, err := db.Prepare(`SELECT user_id, email, first_name, last_name, zip_code,
	                                  is_active, create_dt, update_dt
	                              FROM tinyplannr_api.user
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

func createUserAuthDb(db *sql.DB, u UserCreate) {

	query_str, err := db.Prepare(`INSERT INTO tinyplannr_auth.user
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

	query_str, err := db.Prepare(`INSERT INTO tinyplannr_api.user
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
	createUserAuthDb(db, u)

	lastId := u.ID

	retval, err := getUserDb(db, lastId)
	if err != nil {
		log.Fatal(err)
	}

	return retval, err
}

func getEventById(db *sql.DB, id int64) (*Event, error) {
	var retval Event

	query_str, err := db.Prepare(`SELECT event_id, user_id, title, description,
	                                  location, all_day, start_dt, end_dt, create_dt, end_dt
	                              FROM tinyplannr_api.event
	                              WHERE event_id = $1`)
	if err != nil {
		panic(err)
	}

	err = query_str.QueryRow(id).Scan(&retval.ID, &retval.UserId, &retval.Title, &retval.Description, &retval.Location,
	        &retval.AllDay, &retval.StartDt, &retval.EndDt, &retval.CreateDt, &retval.UpdateDt)
	if err != nil {
		return &retval, err
	}

	return &retval, err
}

func createEventDb(db *sql.DB, e Event) (*Event, error) {
	query_str, err := db.Prepare(`INSERT INTO tinyplannr_api.event VALUES
	                                 (DEFAULT, $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	                              RETURNING event_id`)

	if err != nil {
		panic(err)
	}

	err = query_str.QueryRow(e.UserId, e.Title, e.Description, e.Location, e.AllDay, e.StartDt, e.UpdateDt).Scan(&e.ID)
	if err != nil {
		panic(err)
	}

	lastId := e.ID

	eventData, err := getEventById(db, lastId)
	if err != nil {
		panic(err)
	}

	return eventData, err
}

func loginDb(db *sql.DB, ul UserLogin) (string, error) {
	var hash_pw []byte
	var email string

	password := []byte(ul.Password)

	query_str, err := db.Prepare(`SELECT email, hash_pw FROM tinyplannr_auth.user WHERE email = $1`)
	if err != nil {
		panic(err)
	}

	err = query_str.QueryRow(ul.UserName).Scan(&email, &hash_pw)
	if err != nil {
		panic(err)
	}

	// Compare the hash and PW using the bcrypt library
	// error_str := "Password is incorrect. Please try again."
	err = bcrypt.CompareHashAndPassword(hash_pw, password); if err != nil {
		log.Fatal(err)
	}

	return email, err


}

func createSessionDb(db *sql.DB, sd SessionData) (string, error) {
	var sessionKey string

	// Let's create a session
	session_str, err := db.Prepare(`INSERT INTO tinyplannr_auth.session (session_key, email, update_dt, expire_dt) VALUES
	                                   ($1, $2, CURRENT_TIMESTAMP, $4) RETURNING session_key`)
	if err != nil {
		panic(err)
	}

	err = session_str.QueryRow(sd.SessionId, sd.Username, sd.ExpTime).Scan(&sessionKey)
	if err != nil {
		panic(err)
	}

	return sessionKey, err

}
