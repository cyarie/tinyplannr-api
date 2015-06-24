package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// This struct lets us avoid using global variables all over the place, and instead lets us declare a context we can
// pass into our appHandler struct
type appContext struct {
	db            *sql.DB
	cookieMachine *securecookie.SecureCookie
	handlerResp   int
}

// This struct holds a pointer back to appContext, tells us if the route needs auth or not, gives the handler name
// and holds our extended handler function
type appHandler struct {
	*appContext
	auth_route bool
	route_name string
	h          func(*appContext, http.ResponseWriter, *http.Request)
}

func (fn *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if fn.auth_route == true {
		sid, err := getSessionId(fn.appContext, r)
		if err != nil {
			json.NewEncoder(w).Encode(jsonErr{http.StatusUnauthorized, "Could not validate session. Please try again."})
			log.Println(err)
			return
		}
		sessionCheck, err := validateSessionDb(fn.appContext.db, sid)
		if err != nil {
			json.NewEncoder(w).Encode(jsonErr{http.StatusUnauthorized, "Could not validate session. Please try again."})
			log.Println(err)
			return
		}
		if sessionCheck == true {
			log.Printf("AUTH SUCCESSFUL")
			fn.h(fn.appContext, w, r)
		} else {
			log.Printf("Error: %s", err)
		}
	} else {
		fn.h(fn.appContext, w, r)
	}

	log.Printf(
		"%s\t%s\t%d\t%s\t%s",
		r.Method,
		r.RequestURI,
		fn.appContext.handlerResp,
		fn.route_name,
		time.Since(start),
	)
}

func main() {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sql.Open("postgres", connect_str)

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &appContext{
		db:            db,
		cookieMachine: securecookie.New(cookie_key, cookie_block),
	}

	router := ApiRouter(context)

	context.db.Ping()

	fmt.Println("Connected to the DB... API running!")

	defer context.db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}
