package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type appContext struct {
	db						*sql.DB
	cookieMachine			*securecookie.SecureCookie
	handlerResp				int
}

type appHandler struct{
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request)
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fn.h(fn.appContext, w, r)

	log.Printf(
		"%s\t%s\t%d\t%s",
		r.Method,
		r.RequestURI,
		fn.appContext.handlerResp,
		time.Since(start),
	)
}

func main() {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sql.Open("postgres", connect_str)

	context := &appContext{
		db:					db,
		cookieMachine:		securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
	}

	router := ApiRouter(context)

	context.db.Ping()

	fmt.Println("Connected to the DB...")

	defer context.db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}
