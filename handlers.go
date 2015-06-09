package main

import (
	"net/http"
	"fmt"

	_ "github.com/lib/pq"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME TO GORT")
}