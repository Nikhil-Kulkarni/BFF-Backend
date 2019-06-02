package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route a route
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// Routes a list of routes
type Routes []Route

// NewRouter returns a router with registered routes
func NewRouter(controller *Controller) *mux.Router {
	router := mux.NewRouter()

	var routes = Routes{
		Route{
			"Login",
			"POST",
			"/Login",
			controller.Login,
		},
		Route{
			"FetchScores",
			"POST",
			"/FetchScores",
			controller.FetchScores,
		},
		Route{
			"FetchAllQuestions",
			"POST",
			"/FetchQuestions",
			controller.FetchAllQuestions,
		},
		Route{
			"CreateTest",
			"POST",
			"/CreateTest",
			controller.CreateTest,
		},
		Route{
			"FetchTest",
			"POST",
			"/FetchTest",
			controller.FetchTest,
		},
		Route{
			"SubmitScore",
			"POST",
			"/SubmitScore",
			controller.SubmitScore,
		},
	}

	for _, r := range routes {
		router.
			Methods(r.Method).
			Name(r.Name).
			Path(r.Pattern).
			Handler(r.Handler)
	}

	return router
}
