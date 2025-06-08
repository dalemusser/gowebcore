package asset

import "testing"

func TestPathLookup(t *testing.T) {
	got := Path("app.css")
	if got != "/assets/app.f3c9e2.css" {
		t.Fatalf("unexpected path %s", got)
	}
}
