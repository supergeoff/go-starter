package components

import (
	"fmt"
	"strings"
)

type ButtonProps struct {
	Variant      string // e.g., "default", "destructive", "outline", "secondary", "ghost", "link"
	Size         string // e.g., "default", "sm", "lg", "icon"
	Text         string // The text content of the button
	Href         string // If provided, the button will render as an <a> tag
	Disabled     bool   // If the button should be disabled
	ExtraClasses string // Any additional CSS classes to apply
	Type         string // e.g., "button", "submit", "reset" (defaults to "button" if not an 'a' tag)
}

// GetButtonClasses calculates and returns the combined CSS classes for a button.
func (p ButtonProps) GetButtonClasses() string {
	variant := p.Variant
	if variant == "" {
		variant = "default"
	}

	size := p.Size
	if size == "" {
		size = "default"
	}

	baseClasses := "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity50"

	variantClasses := ""
	switch variant {
	case "default":
		variantClasses = "bg-primary text-primary-foreground shadow hover:bg-primary/90"
	case "destructive":
		variantClasses = "bg-red-500 text-white shadow hover:bg-red-600/90"
	case "outline":
		variantClasses = "border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground"
	case "secondary":
		variantClasses = "bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80"
	case "ghost":
		variantClasses = "hover:bg-accent hover:text-accent-foreground"
	case "link":
		variantClasses = "text-primary underline-offset-4 hover:underline"
	case "success":
		variantClasses = "bg-green-500 text-white shadow hover:bg-green-600/90"
	}

	sizeClasses := ""
	switch size {
	case "default":
		sizeClasses = "h-9 px-4 py-2"
	case "sm":
		sizeClasses = "h-8 rounded-md px-3 text-xs"
	case "lg":
		sizeClasses = "h-10 rounded-md px-8"
	case "icon":
		sizeClasses = "h-9 w-9"
	}

	// Using fmt.Sprintf to combine, then strings.TrimSpace to clean up
	return strings.TrimSpace(
		fmt.Sprintf("%s %s %s %s", baseClasses, variantClasses, sizeClasses, p.ExtraClasses),
	)
}

const ButtonTmplString string = `
{{define "button"}}
    <button
        class="{{.GetButtonClasses}}" {{/* Call the new method */}}>
        {{.Text}}
    </button>
{{end}}
`
