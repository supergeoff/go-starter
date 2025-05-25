package templates

import (
	"log/slog"

	"github.com/supergeoff/go-starter/apps/client/templates/components"
)

// HomePageData defines the structure of data expected by the home template.
type HomePageData struct {
	ButtonData components.ButtonProps
	// Add other fields specific to the home page here
}

const tmplString string = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="min-h-screen flex flex-col items-center justify-center p-8">
    <h1 class="text-4xl font-bold mb-8">Health Check</h1>
    {{template "button" .ButtonData}} {{/* Pass button-specific data to button template */}}
</body>
</html>
`

func init() {
	// Define the components this page template uses
	componentStrings := map[string]string{
		"button": components.ButtonTmplString,
		// Add other components here:
		// "anotherComponent": components.AnotherComponentTmplString,
	}

	// Load the "home" template, providing its own string and the component strings
	LoadTemplate("home", tmplString, componentStrings)
}

// Home prepares the home template for rendering with the given data.
// The data parameter should be of type HomePageData.
// It returns a TemplateRenderer, which has a Render method.
// It panics if the "home" template is not found in the registry,
// which would indicate an issue with the init loading process.
func Home(data interface{}) *TemplateRenderer {
	renderer, err := getRenderer("home", data)
	if err != nil {
		slog.Error("failed to get renderer for home template", "error", err)
		// Consider how to handle this panic in a real app, but for init-time setup, panic is okay.
		panic("Failed to get renderer for home template: " + err.Error())
	}
	return renderer
}
