package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass nil as the third parameter.
	req, err := http.NewRequest("GET", "/testapi", nil)
	assert.NoError(t, err) // Check if an error occurred during request creation.

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ApiHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the Content-Type header.
	assert.Equal(
		t,
		"application/json",
		rr.Header().Get("Content-Type"),
		"handler returned wrong content type",
	)

	// Check the response body is what we expect.
	expectedBody := map[string]string{"message": "check"}
	var actualBody map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err, "failed to unmarshal response body")
	assert.Equal(t, expectedBody, actualBody, "handler returned unexpected body")
}
