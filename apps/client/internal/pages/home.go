package pages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// appelle le backend
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
	err = templates.HomeTempl(status).Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error rendering home: %v\n", err)
	}
}
