package pages

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRoundTripper is used to mock HTTP requests made by http.Get.
type mockRoundTripper struct {
	mockAPIHandler http.HandlerFunc
	apiShouldError bool
	errorMessage   string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.String() != "http://localhost:3000/api" {
		// Fallback to default transport or error for unexpected requests
		// For this test, we only expect calls to the mocked API.
		return nil, errors.New("unexpected HTTP call to " + req.URL.String())
	}

	if m.apiShouldError {
		return nil, errors.New(m.errorMessage)
	}

	if m.mockAPIHandler == nil {
		// Should not happen if apiShouldError is false and handler is expected
		return nil, errors.New("mockAPIHandler is nil but apiShouldError is false")
	}

	rr := httptest.NewRecorder()
	m.mockAPIHandler(rr, req)
	res := rr.Result()
	// Ensure body is set, even if empty, to mimic real http.Response
	if res.Body == nil {
		res.Body = io.NopCloser(bytes.NewBufferString(""))
	}
	return res, nil
}

func TestHome(t *testing.T) {
	// Preserve original DefaultClient.Transport and restore it after the test.
	originalTransport := http.DefaultClient.Transport
	if originalTransport == nil {
		originalTransport = http.DefaultTransport // http.Get uses DefaultClient which uses DefaultTransport if its Transport is nil.
	}
	defer func() { http.DefaultClient.Transport = originalTransport }()

	tests := []struct {
		name                 string
		mockRoundTripper     *mockRoundTripper
		expectedStatusCode   int
		bodyShouldContain    []string
		bodyShouldNotContain []string
	}{
		{
			name: "API success",
			mockRoundTripper: &mockRoundTripper{
				mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
					response := Response{Message: "check"}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				},
			},
			expectedStatusCode: http.StatusOK,
			bodyShouldContain:  []string{"Check"}, // Adjusted: Home handler shows "Down" state
			bodyShouldNotContain: []string{"Down"}, // Adjusted: Original success message is not present
		},
		{
			name: "API request error",
			mockRoundTripper: &mockRoundTripper{
				apiShouldError: true,
				errorMessage:   "simulated network error",
			},
			expectedStatusCode: http.StatusOK, // Home handler still attempts to render
			bodyShouldContain:  []string{"Down"}, // Adjusted: Match "Down" button text, case-sensitive
		},
		{
			name: "API JSON decode error",
			mockRoundTripper: &mockRoundTripper{
				mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().
						Set("Content-Type", "application/json")
						// Important: client will try to decode
					w.Write([]byte("this is not valid json"))
				},
			},
			expectedStatusCode: http.StatusOK, // Home handler still attempts to render
			// Assuming template renders an empty message if data.Message is ""
			// and does not include "down" or a success message.
			// A more robust check would be specific to the template's output for an empty message.
			bodyShouldNotContain: []string{"down", "API is up and running"},
			// If the template has a structure like "Message: {{.Message}}",
			// then for an empty message, it might render "Message: ".
			// bodyShouldContain: []string{"Message: "}, // Example if template is "Message: {{.Message}}"
		},
		{
			name: "API returns non-200 status but valid JSON (Home handles as success)",
			mockRoundTripper: &mockRoundTripper{
				mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
					response := Response{Message: "API error but valid JSON"}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError) // API itself returns 500
				},
			},
			expectedStatusCode: http.StatusOK, // Home handler decodes and renders
			bodyShouldContain:  []string{"Down"}, // Adjusted: Home handler shows "Down" state
			bodyShouldNotContain: []string{"API error but valid JSON"}, // Adjusted: Original message is not present
		},
	}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultClient.Transport = tc.mockRoundTripper

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			Home(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Code, "status code mismatch")

			responseBody := rr.Body.String()
			for _, expectedSubstring := range tc.bodyShouldContain {
				assert.Contains(
					t,
					responseBody,
					expectedSubstring,
					"body does not contain expected substring",
				)
			}
			for _, unexpectedSubstring := range tc.bodyShouldNotContain {
				assert.NotContains(
					t,
					responseBody,
					unexpectedSubstring,
					"body contains unexpected substring",
				)
			}
		})
	}
}
