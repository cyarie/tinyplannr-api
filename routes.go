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
		appHandler{h: Index, auth_route: false},
	},
	Route{
		"User",
		"GET",
		"/user/{userId}",
		appHandler{h: UserIndex, auth_route: true},
	},
	Route{
		"CreateUser",
		"POST",
		"/user/create",
		appHandler{h: CreateUser, auth_route: false},
	},
	Route{
		"CreateEvent",
		"POST",
		"/event/create",
		appHandler{h: CreateEvent, auth_route: true},
	},
	Route{
		"Login",
		"POST",
		"/login",
		appHandler{h: Login, auth_route: false},
	},
}
