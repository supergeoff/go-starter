package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:3000/api")
	status := "down"
	if err == nil && resp.StatusCode == 200 {
		var b map[string]string
		if json.NewDecoder(resp.Body).Decode(&b) == nil && b["message"] == "check" {
			status = "check"
		}
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}

	tmp, err := template.ParseFiles("templates/home.tmpl")
	if err != nil {
		fmt.Printf("Error parsing home: %v\n", err)
	}
	err = tmp.Execute(w, struct{ Status string }{Status: status})
	if err != nil {
		fmt.Printf("Error executing home: %v\n", err)
	}
}
