// Home to all of our HTTP handlers. These lean pretty heavily on the functions in the controllers.go file, which
// contains the functions used to make calls to the database.

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/gorilla/mux"
)

type UserLogin struct {
	UserName		string
	Password		string
}

type LoginResponse struct {
	Email			string	`json:"email"`
}

type SessionData struct {
	SessionId			string
	Username			string
	ExpTime				time.Time
}

func Index(a *appContext, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME TO GORT")
}

func UserIndex(a *appContext, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId int64
	var err error

	// Let's test our sessions by protecting the user index! First, let's grab the session ID from the cookie.
	sid := getSessionId(a, r)

	// Alright, now that we've passed some error checking, let's party
	sessionCheck := validateSessionDb(a.db, sid)

	if sessionCheck == true {
		log.Println("Yay!")
	} else {
		log.Fatal("CHECK FAILED")
	}

	if userId, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		panic(err)
	}

	fmt.Println(userId)
	user, err := getUserDb(a.db, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			a.handlerResp = http.StatusNotFound
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No user found for that ID. Please try again"}); err != nil {
				fmt.Println("faerts")
			}

			return
		}
	}

	if user.ID > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		a.handlerResp = http.StatusOK
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}

		return
	}
}

func CreateUser(a *appContext, w http.ResponseWriter, r *http.Request) {
	var user UserCreate

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		a.handlerResp = 422
		w.WriteHeader(422) // status code for an unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	u, err := createUserDb(a.db, user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	a.handlerResp = http.StatusCreated
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}

func CreateEvent(a *appContext, w http.ResponseWriter, r *http.Request) {
	var event Event

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &event); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		a.handlerResp = 422
		w.WriteHeader(422) // status code for an unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
	}

	e, err := createEventDb(a.db, event)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	a.handlerResp = http.StatusCreated
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Fatal(err)
	}
}

func Login(a *appContext, w http.ResponseWriter, r *http.Request) {
	var ul UserLogin
	var lr LoginResponse
	var sd SessionData

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &ul); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		a.handlerResp = 422
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
	}

	// Let's check that password and make sure it's valid!
	lr.Email, err = loginDb(a.db, ul)
	if err != nil {
		log.Fatal(err)
	}

	// Now that we've passed the login check, let's generate the data we'll fill the cookie with.
	sd.Username = lr.Email
	// Set the session to expire in one month.
	sd.ExpTime = time.Now().UTC().Add(30 * 24 * time.Hour)
	key_str := sd.Username + fmt.Sprint(sd.ExpTime)
	sd.SessionId = generateSessionId(key_str, []byte("as"))

	// Alright, let's write the session to the database
	sk, err := createSessionDb(a.db, sd)
	if err != nil {
		log.Fatal(err)
	}

	// Now that we have a session written to the database, and it has returned a session key/ID for us, let's
	// write that to the cookie and add it to the response
	setSession(a, sk, w)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	a.handlerResp = http.StatusCreated
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(lr); err != nil {
		log.Fatal(err)
	}
}
