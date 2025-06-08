package asset

import (
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

// Handler serves /assets/* from the embedded dist folder.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/assets/")
		if name == "" {
			http.NotFound(w, r)
			return
		}
		data, found := File(name)
		if !found {
			http.NotFound(w, r)
			return
		}

		// set aggressive cache if filename looks fingerprinted
		if strings.ContainsRune(path.Base(name), '.') {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=300")
		}
		http.ServeContent(w, r, name, time.Now(), bytesReader{data})
	})
}

// bytesReader satisfies io.ReadSeeker for ServeContent.
type bytesReader struct{ b []byte }

func (br bytesReader) Read(p []byte) (int, error) { return copy(p, br.b), io.EOF }
func (br bytesReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		return 0, nil
	case io.SeekEnd:
		return int64(len(br.b)), nil
	default:
		return 0, nil
	}
}
