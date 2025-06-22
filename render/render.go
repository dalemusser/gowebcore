// render/render.go
package render

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dalemusser/gowebcore/logger"
)

/*
   ──────────────────────────────────────────────────────────────────────────
      Template loading
   ──────────────────────────────────────────────────────────────────────────
*/

//go:embed templates/*
var Templates embed.FS

// funcs is declared in funcs.go (dict helper, …)
var tpl = template.Must(
	template.
		New("").
		Funcs(funcs). // helper functions
		ParseFS(
			Templates,
			"templates/*.html",             // top-level pages & base layout
			"templates/**/*.html",          // sub-folders
			"templates/_components/*.html", // reusable partials
		),
)

/*
   ──────────────────────────────────────────────────────────────────────────
      Public helpers
   ──────────────────────────────────────────────────────────────────────────
*/

// Render detects HTMX (HX-Request: true) requests and automatically
// chooses between a full HTML page (layout.html) and just the named
// fragment. `name` must correspond to a {{ define }} block.
//
//	render.Render(w, r, "users/table", data)          // fragment
//	render.Render(w, r, "dashboard", data)            // full page
func Render(w http.ResponseWriter, r *http.Request, name string, data any) {
	isHX := r.Header.Get("HX-Request") == "true"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// always inject Now timestamp & make sure data is a map if we need it
	if data == nil {
		data = map[string]any{}
	}
	if m, ok := data.(map[string]any); ok {
		m["Now"] = time.Now().Format(time.RFC3339)
	}

	var err error
	if isHX {
		// HTMX fragment only
		err = tpl.ExecuteTemplate(w, name, data)
	} else {
		// Full page: render layout.html and embed the body template
		err = tpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"Body": template.HTML(name), // name of {{ define }} to include
			"Data": data,
			"Now":  time.Now().Format(time.RFC3339),
		})
	}

	if err != nil {
		logger.Instance().Error("render", "template", name, "err", err)
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
	}
}
