package main

import (
	"github.com/gorilla/mux"
)

func ApiRouter(c *appContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler appHandler

		handler = appHandler{c, route.HandlerFunc.auth_route, route.Name, route.HandlerFunc.h}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(&handler)
	}

	return router
}
