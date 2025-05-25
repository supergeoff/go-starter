package components

// ButtonTmplString exports the raw template string for the button component.
const ButtonTmplString string = `
{{define "button"}}
    <button class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-10 px-4 py-2 {{.ColorClass}} text-white">
        {{.ButtonText}}
    </button>
{{end}}
`
