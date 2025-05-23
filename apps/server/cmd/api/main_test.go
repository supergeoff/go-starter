package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	r := setupRouter()
	foundAPIGet := false
	err := chi.Walk(
		r,
		func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			assert.NotNil(t, handler)
			if method == "GET" && route == "/api" {
				foundAPIGet = true
			}
			return nil
		},
	)
	assert.NoError(t, err, "chi.Walk should not return an error")
	assert.True(t, foundAPIGet, "Expected GET /api route to be registered in the router")
}
