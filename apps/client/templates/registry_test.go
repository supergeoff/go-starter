package templates

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockErrorWriter is an io.Writer that returns an error on Write.
type mockErrorWriter struct{}

func (mew *mockErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mock writer error")
}

// resetGlobalRegistry is a helper function to clean up the global registry for tests.
// This is important because LoadTemplate modifies global state and panics on duplicate names.
func resetGlobalRegistryForTest() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.templates = make(map[string]*template.Template)
}

func TestTemplateRenderer_Render(t *testing.T) {
	// Ensure a clean slate before and after this test suite if other tests modify globalRegistry.
	// However, individual test cases will manage their own state.
	// t.Cleanup(resetGlobalRegistryForTest) // Optional: for overall suite cleanup if needed.

	tests := []struct {
		name           string
		setupRenderer  func(t *testing.T) *TemplateRenderer
		writer         io.Writer
		expectedError  string
		containsError  string // For errors where exact match is hard (like template execution errors)
		expectedOutput string
		cleanup        func() // Specific cleanup for the test case
	}{
		{
			name: "successful render",
			setupRenderer: func(t *testing.T) *TemplateRenderer {
				resetGlobalRegistryForTest() // Clean before this test's setup
				// Pass an empty map for componentTmplStrings as this test doesn't involve components
				LoadTemplate("success_render_test", "Hello {{.Name}}", nil)
				renderer, err := getRenderer(
					"success_render_test",
					map[string]string{"Name": "RenderTest"},
				)
				require.NoError(t, err, "Setup: getRenderer should not fail")
				require.NotNil(t, renderer, "Setup: renderer should not be nil")
				return renderer
			},
			writer:         &bytes.Buffer{},
			expectedOutput: "Hello RenderTest",
			cleanup:        resetGlobalRegistryForTest, // Clean after this test
		},
		{
			name: "render with component",
			setupRenderer: func(t *testing.T) *TemplateRenderer {
				resetGlobalRegistryForTest()
				componentTmpl := `{{define "comp"}}Component: {{.Value}}{{end}}`
				pageTmpl := `Page content. {{template "comp" .}}`
				LoadTemplate("page_with_comp", pageTmpl, map[string]string{"comp": componentTmpl})
				renderer, err := getRenderer(
					"page_with_comp",
					map[string]string{"Value": "TestValue"},
				)
				require.NoError(t, err, "Setup: getRenderer should not fail for component test")
				require.NotNil(t, renderer, "Setup: renderer should not be nil for component test")
				return renderer
			},
			writer:         &bytes.Buffer{},
			expectedOutput: "Page content. Component: TestValue",
			cleanup:        resetGlobalRegistryForTest,
		},
		{
			name: "nil template in renderer",
			setupRenderer: func(t *testing.T) *TemplateRenderer {
				// Directly create a renderer with a nil template to test this specific error path.
				return &TemplateRenderer{template: nil, data: "some data"}
			},
			writer:        &bytes.Buffer{},
			expectedError: "template is not initialized for renderer",
			cleanup:       nil, // No global state modified that needs reset for this specific case
		},
		{
			name: "template execution error due to writer error",
			setupRenderer: func(t *testing.T) *TemplateRenderer {
				resetGlobalRegistryForTest()
				// Pass an empty map for componentTmplStrings as this test doesn't involve components
				LoadTemplate("writer_error_render_test", "Content", nil)
				renderer, err := getRenderer("writer_error_render_test", nil)
				require.NoError(t, err, "Setup: getRenderer should not fail for writer error test")
				require.NotNil(
					t,
					renderer,
					"Setup: renderer should not be nil for writer error test",
				)
				return renderer
			},
			writer:        &mockErrorWriter{},
			containsError: "mock writer error", // This error comes from template.Execute when the writer fails
			cleanup:       resetGlobalRegistryForTest,
		},
		{
			name: "template execution error due to bad template action",
			setupRenderer: func(t *testing.T) *TemplateRenderer {
				// This test does not use the global registry, it creates a template directly.
				tmpl, err := template.New("bad_action_render_test").
					Parse("Hello {{.NonExistentField}}")
				require.NoError(
					t,
					err,
					"Setup: parsing template for bad action test should not fail",
				)
				// Data struct does not have NonExistentField
				return &TemplateRenderer{template: tmpl, data: struct{ Name string }{"Test"}}
			},
			writer: &bytes.Buffer{},
			// The exact error message can be complex and might include line numbers,
			// so we check for a key part of the error.
			containsError: `executing "bad_action_render_test" at <.NonExistentField>: can't evaluate field NonExistentField`,
			cleanup:       nil, // No global state modified
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cleanup != nil {
				// Defer cleanup to run after the test case finishes
				defer tc.cleanup()
			}

			renderer := tc.setupRenderer(t)
			require.NotNil(t, renderer, "Renderer should not be nil after setup")

			err := renderer.Render(tc.writer)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError, "Error message mismatch")
			} else if tc.containsError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.containsError, "Error message should contain expected substring")
			} else {
				assert.NoError(t, err, "Render should not produce an error")
				if buf, ok := tc.writer.(*bytes.Buffer); ok {
					assert.Equal(t, tc.expectedOutput, buf.String(), "Rendered output mismatch")
				}
			}
		})
	}
}
