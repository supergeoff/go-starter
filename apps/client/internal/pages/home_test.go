package pages

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock Response struct to match the one in home.go
type MockResponse struct {
	Message string `json:"message"`
}

// mockRoundTripper allows controlling HTTP responses for testing.
type mockRoundTripper struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (mrt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if mrt.RoundTripFunc != nil {
		return mrt.RoundTripFunc(req)
	}
	return nil, errors.New("RoundTrip not implemented")
}

func TestHome(t *testing.T) {
	tests := []struct {
		name                string
		apiStatusCode       int
		apiResponseBody     interface{}
		apiError            error // To simulate http.Get errors
		expectedStatus      int
		expectedButtonClass string
		expectedButtonText  string
	}{
		{
			name:                "successful API call - check message",
			apiStatusCode:       http.StatusOK,
			apiResponseBody:     MockResponse{Message: "check"},
			apiError:            nil,
			expectedStatus:      http.StatusOK,
			expectedButtonClass: "bg-green-500",
			expectedButtonText:  "OK",
		},
		{
			name:                "successful API call - other message",
			apiStatusCode:       http.StatusOK,
			apiResponseBody:     MockResponse{Message: "hello"},
			apiError:            nil,
			expectedStatus:      http.StatusOK,
			expectedButtonClass: "bg-red-600",
			expectedButtonText:  "Down",
		},
		{
			name:                "API call failed",
			apiStatusCode:       0, // Status code is irrelevant if the call fails
			apiResponseBody:     nil,
			apiError:            errors.New("mock http get error"), // Simulate an error
			expectedStatus:      http.StatusOK,                     // Home function still renders template on API error
			expectedButtonClass: "bg-red-600",
			expectedButtonText:  "Down",
		},
		{
			name:                "JSON decoding failed",
			apiStatusCode:       http.StatusOK,
			apiResponseBody:     `{"message": 123}`, // Invalid JSON for the struct
			apiError:            nil,
			expectedStatus:      http.StatusOK, // Home function still renders template on JSON error
			expectedButtonClass: "bg-red-600",
			expectedButtonText:  "Down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the HTTP client's RoundTripper
			originalTransport := http.DefaultClient.Transport
			http.DefaultClient.Transport = &mockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					if tt.apiError != nil {
						return nil, tt.apiError
					}

					// Create a mock response
					recorder := httptest.NewRecorder()
					recorder.WriteHeader(tt.apiStatusCode)
					if tt.apiResponseBody != nil {
						// Handle string body for JSON decoding failure test
						if bodyStr, ok := tt.apiResponseBody.(string); ok {
							_, err := recorder.WriteString(bodyStr)
							if err != nil {
								slog.Error("Error writting string", "error", err)
							}
						} else {
							err := json.NewEncoder(recorder).Encode(tt.apiResponseBody)
							if err != nil {
								slog.Error("Error encoding json", "error", err)
							}
						}
					}

					return recorder.Result(), nil
				},
			}
			defer func() { http.DefaultClient.Transport = originalTransport }() // Restore original transport

			// Create a ResponseRecorder and a dummy Request
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			// Call the Home function
			Home(rr, req)

			// Assert the status code and body
			assert.Equal(t, tt.expectedStatus, rr.Code, "Handler returned wrong status code")

			// Assert on specific elements instead of the full HTML
			bodyString := rr.Body.String()
			assert.Contains(t, bodyString, "<button", "Response body should contain a button")
			assert.Contains(
				t,
				bodyString,
				tt.expectedButtonClass,
				"Button should have the expected class",
			)
			assert.Contains(
				t,
				bodyString,
				tt.expectedButtonText,
				"Button should have the expected text",
			)
		})
	}
}
