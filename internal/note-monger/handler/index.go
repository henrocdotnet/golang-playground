package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/henrocdotnet/golang-playground/internal/note-monger/renderer"
	"github.com/henrocdotnet/golang-playground/internal/note-monger/textnote"
)

func Index(c echo.Context) error {
	// Setup template data.
	data := make(map[string]interface{})
	data["list"] = []textnote.TextNote{{ Title: "One", Body: "One"}, { Title: "Two", Body: "Two" }}

	// Render templates.
	out, err := renderer.RenderWithLayout(c, "index", "layout.html", "index.html", data)
	if err != nil {
		log.Printf("Could not render page index: %v", err)
		return err
	}

	return c.HTML(http.StatusOK, out.String())
}