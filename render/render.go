package render

import (
	"embed"
	"html/template"
	"net/http"
	"time"
)

//go:embed templates/*
var Templates embed.FS

var tpl = template.Must(
	template.ParseFS(
		Templates,
		"templates/*.html",     // top-level templates
		"templates/**/*.html",  // nested templates
	),
)

// Render chooses full page vs. HTMX fragment automatically.
func Render(w http.ResponseWriter, r *http.Request, name string, data any) {
	isHX := r.Header.Get("HX-Request") == "true"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if data == nil {
		data = map[string]any{}
	}
	if m, ok := data.(map[string]any); ok {
		m["Now"] = time.Now().Format(time.RFC3339)
	}

	_ = tpl.ExecuteTemplate(w, name, data)

	// For HTMX fragments the named template is the fragment itself.
	// For full pages the template usually includes {{ template "base" . }}.
	if !isHX {
		// nothing extra to do – ExecuteTemplate already rendered base layout.
	}
}
