package components

// ButtonTmplString exports the raw template string for the button component.
const ButtonTmplString string = `
{{define "button"}}
    {{$buttonText := "Down"}}
    {{$colorClass := "bg-red-600"}}
    {{if eq .Message "check"}}
        {{$buttonText = "OK"}}
        {{$colorClass = "bg-green-500"}}
    {{end}}
    <button class="{{$colorClass}} text-white px-4 py-2 rounded text-3xl">
        {{$buttonText}}
    </button>
{{end}}
`
