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

// LoadTemplate parses a page template string and any provided component template strings,
// storing the resulting composite template with the given name in the global registry.
// It panics if parsing fails or if the template name is already registered.
func LoadTemplate(name string, pageTmplString string, componentTmplStrings map[string]string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if _, ok := globalRegistry.templates[name]; ok {
		slog.Error("page template with that name already exists", "template", name)
		panic("Error: page template with that name already exists: " + name)
	}

	// Create a new template. This will be the container for the page and its components.
	tmpl := template.New(name)

	// Parse all provided component template strings into this page's template set.
	for componentName, componentStr := range componentTmplStrings {
		// The Parse method adds the definitions from componentStr to tmpl.
		// If componentStr contains {{define "compName"}}, "compName" becomes available.
		_, err := tmpl.Parse(componentStr)
		if err != nil {
			slog.Error("failed to parse component template into page template",
				"page", name, "component", componentName, "error", err)
			panic(
				"Error: failed to parse component template '" + componentName + "' for page '" + name + "': " + err.Error(),
			)
		}
		slog.Info(
			"parsed component into page template",
			"page",
			name,
			"component",
			componentName,
		)
	}

	// Now parse the main page template string itself into the same template set.
	_, err := tmpl.Parse(pageTmplString)
	if err != nil {
		slog.Error("failed to parse main page template string", "template", name, "error", err)
		panic("Error: failed to parse main page template '" + name + "': " + err.Error())
	}

	globalRegistry.templates[name] = tmpl
	slog.Info("page template loaded with components", "template", name)
}

// getRenderer retrieves a parsed template by name and prepares it for rendering with the given data.
// This is an internal helper function.
func getRenderer(name string, data interface{}) (*TemplateRenderer, error) {
	globalRegistry.mu.RLock()
	tmpl, ok := globalRegistry.templates[name]
	globalRegistry.mu.RUnlock()

	if !ok {
		slog.Error("template not found in registry", "template", name)
		return nil, errors.New("error: template not found in registry: " + name)
	}
	return &TemplateRenderer{template: tmpl, data: data}, nil
}
