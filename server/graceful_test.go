package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/dalemusser/gowebcore/logger"
)

func TestGraceful(t *testing.T) {
	logger.Init("info") // ensure global logger is non-nil

	srv := &http.Server{Addr: ":0", Handler: http.NewServeMux()}
	l, _ := net.Listen("tcp", srv.Addr)
	go srv.Serve(l)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	if err := Graceful(ctx, srv); err != nil {
		t.Fatalf("graceful shutdown failed: %v", err)
	}
}
