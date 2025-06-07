package server

import (
	"encoding/json"
	"net/http"

	"github.com/dalemusser/gowebcore/version"
)

func DefaultHealthHandler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"commit": version.Commit,
		"build":  version.Build,
		"go":     version.GoVersion,
	})
}
