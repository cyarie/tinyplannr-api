package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME TO GORT")
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId int64
	var err error

	if userId, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		panic(err)
	}

	fmt.Println(userId)
	user, err := getUserById(db, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No user found for that ID. Please try again"}); err != nil {
				fmt.Println("faerts")
			}

			return
		}
	}

	if user.ID > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}

		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreate

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // status code for an unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	u, err := createUserDb(db, user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &event); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // status code for an unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
	}

	e, err := createEventDb(db, event)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Fatal(err)
	}
}
