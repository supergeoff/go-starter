package pages

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/templates"
)

type ButtonData struct {
	ButtonText string
	ColorClass string
}

type HomePageData struct {
	Message    string
	ButtonData ButtonData
}

func Home(w http.ResponseWriter, r *http.Request) {
	var apiResponse struct {
		Message string `json:"message"`
	}
	var pageData HomePageData

	resp, err := http.Get("http://localhost:3000/api")
	if err != nil {
		slog.Error("Error making GET request", "error", err)
		apiResponse.Message = "down"
	} else {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			slog.Error("Failed to decode JSON", "error", err)
		}
	}

	pageData.Message = apiResponse.Message

	// Logic for button
	if pageData.Message == "check" {
		pageData.ButtonData.ButtonText = "OK"
		pageData.ButtonData.ColorClass = "bg-green-500"
	} else {
		pageData.ButtonData.ButtonText = "Down"
		pageData.ButtonData.ColorClass = "bg-red-600"
	}

	err = templates.Home(pageData).Render(w)
	if err != nil {
		slog.Error("Error rendering template", "error", err)
	}
}
