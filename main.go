package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type appContext struct {
	db						*sql.DB
	cookieMachine			*securecookie.SecureCookie
}

type appHandler struct{
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request)
}


func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn.h(fn.appContext, w, r)
}

func main() {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sql.Open("postgres", connect_str)
	context := &appContext{
		db:					db,
		cookieMachine:		securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
	}
	router := ApiRouter()

	context.db.Ping()

	fmt.Println("Connected to the DB...")

	defer context.db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}
