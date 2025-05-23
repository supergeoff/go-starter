package pages

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/templates"
)

type Response struct {
	Message string `json:"message"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	var data Response

	resp, err := http.Get("http://localhost:3000/api")
	if err != nil {
		slog.Error("Error making GET request", "error", err)
		data.Message = "down"
	} else {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			slog.Error("Failed to decode JSON", "error", err)
		}
	}

	err = templates.Home(data).Render(w)
	if err != nil {
		slog.Error("Error rendering template", "error", err)
	}
}
