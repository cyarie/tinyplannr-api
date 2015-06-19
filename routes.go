package main

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appHandler
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		appHandler{h: Index},
	},
	Route{
		"User",
		"GET",
		"/user/{userId}",
		appHandler{h: UserIndex},
	},
	Route{
		"CreateUser",
		"POST",
		"/user/create",
		appHandler{h: CreateUser},
	},
	Route{
		"CreateEvent",
		"POST",
		"/event/create",
		appHandler{h: CreateEvent},
	},
	Route{
		"Login",
		"POST",
		"/login",
		appHandler{h: Login},
	},
}
