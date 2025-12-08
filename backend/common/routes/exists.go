package routes

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	Path   string
	Method string
}

func RouteExists(t *testing.T, routes chi.Router, apiVersion int, route Route) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if fmt.Sprintf("/v%d%s", apiVersion, route.Path) == foundRoute && route.Method == method {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s with method %s in registered routes", route.Path, route.Method)
	}
}
