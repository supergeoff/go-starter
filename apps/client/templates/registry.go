package templates

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"sync"
)

// TemplateRenderer is a struct that holds a specific template and data for rendering.
type TemplateRenderer struct {
	template *template.Template
	data     interface{}
}

// Render executes the template with the associated data and writes to w.
// It returns an error if the template execution fails.
func (tr *TemplateRenderer) Render(w io.Writer) error {
	if tr.template == nil {
		// This should ideally not be reached if template loading and retrieval are correct.
		slog.Error("template is not initialized for renderer")
		return errors.New("template is not initialized for renderer")
	}
	return tr.template.Execute(w, tr.data)
}

// registry manages named templates.
// It's an unexported struct as we'll use a global instance.
type registry struct {
	mu        sync.RWMutex
	templates map[string]*template.Template
}

// globalRegistry is the single, global instance of our template registry.
var globalRegistry = &registry{
	templates: make(map[string]*template.Template),
}

// LoadTemplate parses a template string and stores it with the given name in the global registry.
// It panics if parsing fails or if the template name is already registered,
// as these are considered critical setup errors during application initialization.
func LoadTemplate(name string, tmplString string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if _, ok := globalRegistry.templates[name]; ok {
		slog.Error("template with name that name alreay exists", "template", name)
		panic("Error: template with that name already exists")
	}

	tmpl, err := template.New(name).Parse(tmplString)
	if err != nil {
		slog.Error("impossible to parse template", "template", name, "error", err)
		panic("Error: failed to parse template")

	}
	globalRegistry.templates[name] = tmpl
}

// getRenderer retrieves a parsed template by name and prepares it for rendering with the given data.
// This is an internal helper function.
func getRenderer(name string, data interface{}) (*TemplateRenderer, error) {
	globalRegistry.mu.RLock()
	tmpl, ok := globalRegistry.templates[name]
	globalRegistry.mu.RUnlock()

	if !ok {
		slog.Error("template not found in registry", "template", name)
		return nil, errors.New("error: template not found in registry")
	}
	return &TemplateRenderer{template: tmpl, data: data}, nil
}
