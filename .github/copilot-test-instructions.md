---
mode: 'chat'
description: 'Generate test for a file'
---

# Prompt for Generating Go Unit Tests

Your goal is to generate comprehensive unit tests for the Go source file provided.

**File to Test:**
- Please specify the file you want to write tests for using the `#file:` directive (e.g., `#file:internal/service/user_service.go`).

**Testing Guidelines:**

1.  **Test File Naming:** Create a test file named `[original_filename]_test.go` in the same package as the source file.
2.  **Package:** Ensure the test file is in the `[original_package]_test` package if you need to test exported identifiers and avoid import cycles, or in the same `[original_package]` if testing unexported identifiers is necessary and safe. Prefer `_test` packages for black-box testing.
3.  **Test Functions:**
    *   Name test functions `Test[FunctionName]`.
    *   For each public function/method in the source file, generate at least one test function.
4.  **Testify Library:**
    *   Utilize the `testify/assert` package for non-fatal assertions (e.g., `assert.NoError(t, err)`, `assert.Equal(t, expected, actual)`).
    *   Utilize the `testify/require` package for fatal assertions where the test cannot proceed if the assertion fails (e.g., `require.NoError(t, err)` before further operations that depend on no error).
5.  **Table-Driven Tests:**
    *   For functions with multiple distinct inputs or scenarios, prefer table-driven tests.
    *   Define a struct for test cases (e.g., `name string, input an_input_type, want an_output_type,wantErr bool`).
    *   Iterate over the table, running each test case as a sub-test using `t.Run()`.
6.  **Coverage:**
    *   Aim for good test coverage, including:
        *   Happy paths (valid inputs, expected outputs).
        *   Edge cases (e.g., nil inputs, empty slices/maps, zero values, very large/small numbers).
        *   Error conditions (functions that return errors should be tested to ensure errors are returned correctly).
7.  **Mocking:**
    *   If the function under test has external dependencies (e.g., database, external API calls, file system interactions), suggest using mocks.
    *   You can use `testify/mock` or standard Go interfaces with mock implementations.
    *   Focus on mocking the *behavior* of dependencies. Consider approaches for setting up and tearing down test environments or mocks if needed.
8.  **HTTP Handlers (if applicable):**
    *   If testing `http.HandlerFunc` or `chi` router handlers, use the `net/http/httptest` package to create mock requests and record responses.
    *   Check status codes, response bodies, and headers.
9.  **Clarity and Readability:**
    *   Generated tests should be clear, concise, and easy to understand.
    *   Add brief comments if the test logic is complex or non-obvious.

**Example Structure (for a hypothetical `calculator.go`):**

```go
// In calculator_test.go
package calculator_test // or package calculator

import (
	"testing"
	"your_module_path/calculator" // Adjust import path
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a       int
		b       int
		want    int
		wantErr bool // If Add could return an error
	}{
		{
			name: "positive numbers",
			a:    2,
			b:    3,
			want: 5,
		},
		{
			name: "negative numbers",
			a:    -2,
			b:    -3,
			want: -5,
		},
		// Add more test cases for edge cases and errors
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Add(tt.a, tt.b) // Assuming Add might return an error
			if tt.wantErr {
				require.Error(t, err)
				// assert.Equal(t, expectedErrorType, err) // Optional: check error type/message
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
```

10. **Provide the following information if not clear from the #file context:**
    * Any specific scenarios or edge cases you want to prioritize for testing.
    * Details about complex data structures or dependencies that might require specific setup or mocking.