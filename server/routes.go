package server

import (
	"fmt"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var serverRoutes = routes{
	route{
		"Index",
		"GET",
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello world")
		},
	},
}
