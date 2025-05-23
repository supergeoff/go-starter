package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	r := setupRouter()
	require.NotNil(t, r, "setupRouter() should return a non-nil chi.Mux router")

	var foundAPIGet bool

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
			if method == "GET" && route == "/api" {
				foundAPIGet = true
			}
			return nil
		},
	)
	assert.NoError(t, err, "chi.Walk should not return an error")
	assert.True(t, foundAPIGet, "Expected GET /api route to be registered in the router")
}
