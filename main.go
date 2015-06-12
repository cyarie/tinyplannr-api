package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

var (
	db *sql.DB
)

func main() {
	router := ApiRouter()

	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ = sql.Open("postgres", connect_str)

	db.Ping()

	fmt.Println("Connected to the DB...")

	defer db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}