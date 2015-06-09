package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"database/sql"
)

func main() {
	router := ApiRouter()

	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, err := sql.Open("postgres", connect_str)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	userTest, err := getUserById(db, 1)
	if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(os.Stdout, "User: %v\n", userTest.Email)
	}


	log.Fatal(http.ListenAndServe(":8080", router))
}
