package server

import "net/http"

type route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var serverRoutes = routes{
	route{
		"GraphQL endpoint",
		[]string{"GET", "POST"},
		"/graphql",
		query,
	},
}
