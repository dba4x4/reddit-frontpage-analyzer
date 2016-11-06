package server

import (
	"net/http"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/server/resource"
)

type route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var serverRoutes = routes{
	route{
		"Get posts with tags filtered by date",
		[]string{"GET"},
		"/api/v1/posts",
		resource.GetPosts,
	},
}
