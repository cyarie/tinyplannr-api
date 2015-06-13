package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"User",
		"GET",
		"/user/{userId}",
		UserIndex,
	},
	Route{
		"CreateUser",
		"POST",
		"/user/create",
		CreateUser,
	},
	Route{
		"CreateEvent",
		"POST",
		"/event/create",
		CreateEvent,
	},
}
