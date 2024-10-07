package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello To The MusicService!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GetSongs",
		"POST",
		"/data",
		GetSongs,
	},

	Route{
		"GetSongText",
		"POST",
		"/song",
		GetSongText,
	},

	Route{
		"Change",
		"POST",
		"/change",
		Change,
	},

	Route{
		"Delete",
		"POST",
		"/delete",
		Delete,
	},

	Route{
		"Add",
		"POST",
		"/add",
		Add,
	},
}
