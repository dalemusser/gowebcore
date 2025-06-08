package asset

import (
	"embed"
	"encoding/json"
	"path"
)

//go:embed dist/** manifest.json
var embedded embed.FS

// manifest maps logical -> hashed filenames.
var manifest map[string]string

func init() {
	data, _ := embedded.ReadFile("manifest.json")
	_ = json.Unmarshal(data, &manifest)
}

// Path returns "/assets/<hashed>" given logical name "app.css".
func Path(logical string) string {
	if hashed, ok := manifest[logical]; ok {
		return "/assets/" + hashed
	}
	// fallback to logical (useful during dev without hashing)
	return "/assets/" + logical
}

// File returns the bytes and a bool (true if found).
func File(name string) ([]byte, bool) {
	b, err := embedded.ReadFile(path.Join("dist", name))
	return b, err == nil
}