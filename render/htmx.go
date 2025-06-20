package render

import (
	"net/http"

	"github.com/dalemusser/gowebcore/logger"
)

// isHX reports true when the request comes from HTMX.
func isHX(r *http.Request) bool { return r.Header.Get("HX-Request") == "true" }

// Page renders a full document unless the request comes from HTMX.
//
//   - Non-HTMX  → executes the child template; its {{ block "..." }} definitions
//     will be merged into base.html automatically.
//   - HTMX      → falls back to Fragment (no layout).
func Page(w http.ResponseWriter, r *http.Request, child string, data any) {
	if isHX(r) {
		Fragment(w, r, child, data)
		return
	}
	// For traditional requests we execute the child template itself.
	// Because tpl.ParseFS loaded base.html first, Go’s template engine
	// injects the child’s {{ block "title" }}, {{ block "content" }}, etc.
	// into base.html and returns the complete HTML page.
	if err := tpl.ExecuteTemplate(w, child, data); err != nil {
		logger.Instance().Error("render page", "err", err)
	}
}

// Fragment renders a single template without layout.
func Fragment(w http.ResponseWriter, r *http.Request, name string, data any) {
	if err := tpl.ExecuteTemplate(w, name, data); err != nil {
		logger.Instance().Error("render fragment", "err", err)
	}
}

// HXRedirect sends HX-Redirect for HTMX or a normal 303.
func HXRedirect(w http.ResponseWriter, r *http.Request, url string) {
	if isHX(r) {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// HXTrigger fires an HTMX client event.
func HXTrigger(w http.ResponseWriter, event string) {
	w.Header().Set("HX-Trigger", event)
}
