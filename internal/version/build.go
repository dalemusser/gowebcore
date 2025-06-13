package version

import "runtime"

var (
	Commit     = "devel"
	BuildDate  = "unknown"
	GoVersion  = runtime.Version()
)
