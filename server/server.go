package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// New returns a ready-to-start *http.Server.
// - cfg provides ports, domain, TLS flag.
// - routes is your service's router; /health is added automatically.
func New(cfg config.Base, routes chi.Router) *http.Server {
	// default middleware stack
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP)
	r.Use(middleware.Recoverer, middleware.Compress(5))
	r.Use(logger.ChiLogger)

	// mount caller routes and health
	r.Mount("/", routes)
	r.Get("/health", DefaultHealthHandler)

	addr := net.JoinHostPort("", pickPort(cfg.HTTPPort, 8080))

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// HTTPS with Let’s Encrypt (if enabled)
	if cfg.EnableTLS && cfg.Domain != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(".cert-cache"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
		}
		go http.ListenAndServe(":http", m.HTTPHandler(nil)) // ACME challenge

		srv.Addr = net.JoinHostPort("", pickPort(cfg.HTTPSPort, 8443))
		srv.TLSConfig = &tls.Config{
			GetCertificate: m.GetCertificate,
			NextProtos:     []string{acme.ALPNProto, "h2", "http/1.1"},
		}
	}
	return srv
}

// Graceful blocks until ctx.Done() or SIGINT/SIGTERM, then shuts the server down.
func Graceful(ctx context.Context, srv *http.Server) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-stop:
	}
	logger.Instance().Info("shutting down")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(c)
}

// pickPort returns cfgPort if non-zero, else the fallback.
func pickPort(cfgPort, fallback int) string {
	if cfgPort == 0 {
		return fmt.Sprintf("%d", fallback)
	}
	return fmt.Sprintf("%d", cfgPort)
}
