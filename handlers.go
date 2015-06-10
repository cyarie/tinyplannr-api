package main

import (
	"net/http"
	"fmt"
	"database/sql"
	"strconv"
	"encoding/json"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME TO GORT")
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId int
	var err error

	if userId, err = strconv.Atoi(vars["userId"]); err != nil {
		panic(err)
	}

	fmt.Println(userId)
	user, err := getUserById(db, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No user found for that ID. Please try again"});
			err != nil {
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