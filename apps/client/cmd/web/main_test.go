package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetupRouter_WebAppRoutes tests the router setup for web application routes
// (GET / and /static/*) based on the structure of the `setupRouter` function
// provided in the `codeToTest` section.
func TestSetupRouter_WebAppRoutes(t *testing.T) {
	// This test assumes that a setupRouter() function, matching the provided
	// codeToTest structure, exists in the current 'main' package.
	r := setupRouter()
	require.NotNil(t, r, "setupRouter() should return a non-nil chi.Mux router")

	var (
		foundIndexGet      bool
		foundStaticRoute   bool
		staticRoutePattern = "/static/*"
	)

	err := chi.Walk(
		r,
		func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			require.NotNil(
				t,
				handler,
				"Handler for method %s, route %s should not be nil",
				method,
				route,
			)

			if method == http.MethodGet && route == "/" {
				foundIndexGet = true
			}

			if route == staticRoutePattern {
				// This confirms that a handler is registered for the "/static/*" pattern.
				// chi.Router.Handle registers for all methods. chi.Walk might list this route
				// multiple times for different specific methods or with a generic method indicator.
				// Confirming the pattern's presence is the key aspect.
				foundStaticRoute = true
			}
			return nil
		},
	)
	require.NoError(t, err, "chi.Walk should not return an error during router traversal")
	assert.True(t, foundIndexGet, "Expected GET / route to be registered")
	assert.True(t, foundStaticRoute, "Expected "+staticRoutePattern+" route to be registered")
}
