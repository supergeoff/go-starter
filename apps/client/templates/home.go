package templates

import (
	"log"
)

const tmplString string = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Home Page</title>
    <link rel="stylesheet" href="/static/css/global.css">
</head>
<body class="p-8">
    {{if eq .Message "check"}}
        <button class="bg-green-500 text-white px-4 py-2 rounded text-3xl">
            OK
        </button>
    {{else}}
        <button class="bg-red-600 text-white px-4 py-2 rounded text-3xl">
            Down
        </button>
    {{end}}
</body>
</html>
`

func init() {
	// LoadTemplate will panic if there's an error (e.g., parsing error, duplicate name),
	// which is appropriate for an init function.
	LoadTemplate("home", tmplString)
}

// Home prepares the home template for rendering with the given data.
// It returns a TemplateRenderer, which has a Render method.
// It panics if the "home" template is not found in the registry,
// which would indicate an issue with the init loading process.
func Home(data interface{}) *TemplateRenderer {
	renderer, err := getRenderer("home", data)
	if err != nil {
		// This panic indicates a programming error, e.g., the template name "home"
		// used here doesn't match the name used in LoadTemplate during init,
		// or LoadTemplate was not called for "home".
		log.Panicf(
			"Failed to get renderer for 'home': %v. Ensure template was loaded correctly during init.",
			err,
		)
	}
	return renderer
}
