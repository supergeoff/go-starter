package pages

import (
	"encoding/json"
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/templates"
)

type indexData struct {
	Status string
}

func Home(w http.ResponseWriter, r *http.Request) {
	// appelle le backend
	resp, err := http.Get("http://localhost:3000/api")
	status := "down"
	if err == nil && resp.StatusCode == 200 {
		var b map[string]string
		if json.NewDecoder(resp.Body).Decode(&b) == nil && b["message"] == "check" {
			status = "check"
		}
		resp.Body.Close()
	}
	templates.HomeTempl(status).Render(r.Context(), w)
}
