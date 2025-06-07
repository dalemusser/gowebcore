package version

import "runtime"

// These variables are meant to be set at build-time via -ldflags.
// Example: go build -ldflags "-X github.com/dalemusser/gowebcore/version.Build=1.2.3 -X github.com/dalemusser/gowebcore/version.Commit=abcdef0 -X github.com/dalemusser/gowebcore/version.Date=2025-06-07T12:00:00Z"
var (
	Build     = "dev"                      // semantic version or tag
	Commit    = "none"                     // git SHA
	Date      = "unknown"                  // build date
	GoVersion = runtime.Version()          // runtime version captured at init
)
