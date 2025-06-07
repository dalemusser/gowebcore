package render

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderFragmentHX(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/fragment", nil)
	r.Header.Set("HX-Request", "true")

	Render(w, r, "fragment.html", map[string]any{"Now": "test"})

	if !strings.Contains(w.Body.String(), "Loaded via HTMX") {
		t.Fatalf("expected fragment content, got %s", w.Body.String())
	}
}
