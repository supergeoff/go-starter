package pages

import (
	"encoding/json"
	"errors"
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
		name            string
		apiStatusCode   int
		apiResponseBody interface{}
		apiError        error // To simulate http.Get errors
		expectedOutput  string
		expectedStatus  int
	}{
		{
			name:            "successful API call - check message",
			apiStatusCode:   http.StatusOK,
			apiResponseBody: MockResponse{Message: "check"},
			apiError:        nil,
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="p-8">
    
        <button class="bg-green-500 text-white px-4 py-2 rounded text-3xl">
            OK
        </button>
    
</body>
</html>
`,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "successful API call - other message",
			apiStatusCode:   http.StatusOK,
			apiResponseBody: MockResponse{Message: "hello"},
			apiError:        nil,
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="p-8">
    
        <button class="bg-red-600 text-white px-4 py-2 rounded text-3xl">
            Down
        </button>
    
</body>
</html>
`,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "API call failed",
			apiStatusCode:   0, // Status code is irrelevant if the call fails
			apiResponseBody: nil,
			apiError:        errors.New("mock http get error"), // Simulate an error
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="p-8">
    
        <button class="bg-red-600 text-white px-4 py-2 rounded text-3xl">
            Down
        </button>
    
</body>
</html>
`,
			expectedStatus: http.StatusOK, // Home function still renders template on API error
		},
		{
			name:            "JSON decoding failed",
			apiStatusCode:   http.StatusOK,
			apiResponseBody: `{"message": 123}`, // Invalid JSON for the struct
			apiError:        nil,
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="p-8">
    
        <button class="bg-red-600 text-white px-4 py-2 rounded text-3xl">
            Down
        </button>
    
</body>
</html>
`,
			expectedStatus: http.StatusOK, // Home function still renders template on JSON error
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
							recorder.WriteString(bodyStr)
						} else {
							json.NewEncoder(recorder).Encode(tt.apiResponseBody)
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
			assert.Equal(t, tt.expectedOutput, rr.Body.String(), "Handler returned unexpected body")
		})
	}
}
