package pages

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/templates"
	"github.com/supergeoff/go-starter/apps/client/templates/components" // Import components for ButtonProps
)

func Home(w http.ResponseWriter, r *http.Request) {
	var apiResponse struct {
		Message string `json:"message"`
	}

	var pageData templates.HomePageData

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

	// Initialize ButtonData with default values
	buttonProps := components.ButtonProps{
		Variant: "default", // Set a default variant
		Size:    "default", // Set a default size
	}

	// Logic for button based on API response
	if apiResponse.Message == "check" {
		buttonProps.Text = "OK"
		buttonProps.Variant = "success"
	} else {
		buttonProps.Text = "Down"
		buttonProps.Variant = "destructive"
	}

	pageData.ButtonData = buttonProps // Assign the prepared buttonProps

	err = templates.Home(pageData).Render(w)
	if err != nil {
		slog.Error("Error rendering template", "error", err)
	}
}
