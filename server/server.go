// server/server.go
package server

import (
	"context"
	"crypto/tls"
	"errors"
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

/* ---------------------------------------------------------------------- */
/*                             Construction                               */
/* ---------------------------------------------------------------------- */

// New returns a *http.Server* pre-configured for
//  1. static-cert TLS   (highest precedence)
//  2. Let’s Encrypt TLS (if enabled)
//  3. plain HTTP        (fallback).
func New(cfg config.Base, routes chi.Router) *http.Server {
	// default middleware stack
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP)
	r.Use(middleware.Recoverer, middleware.Compress(5))
	r.Use(logger.ChiLogger)

	// mount caller routes + health endpoint
	r.Mount("/", routes)
	r.Get("/health", DefaultHealthHandler)

	srv := &http.Server{
		Addr:              net.JoinHostPort("", pickPort(cfg.HTTPPort, 8080)),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	/* ---------- 1 — Static certificate? ---------- */
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			logger.Error("TLS cert load failed", "err", err)
		} else {
			srv.Addr = net.JoinHostPort("", pickPort(cfg.HTTPSPort, 8443))
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
				NextProtos:   []string{"h2", "http/1.1"},
			}
			return srv
		}
	}

	/* ---------- 2 — Let’s Encrypt autocert ---------- */
	if cfg.EnableTLS && cfg.Domain != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(".cert-cache"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
		}
		go http.ListenAndServe(":http", m.HTTPHandler(nil)) // ACME HTTP-01

		srv.Addr = net.JoinHostPort("", pickPort(cfg.HTTPSPort, 8443))
		srv.TLSConfig = &tls.Config{
			GetCertificate: m.GetCertificate,
			NextProtos:     []string{acme.ALPNProto, "h2", "http/1.1"},
		}
	}

	return srv // plain HTTP if we got here
}

/* ---------------------------------------------------------------------- */
/*                Graceful shutdown + Serve helper                        */
/* ---------------------------------------------------------------------- */

// Graceful blocks until ctx.Done() or SIGINT/SIGTERM, then shuts the server.
func Graceful(ctx context.Context, srv *http.Server) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-stop:
	}
	logger.Info("shutting down")

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(c)
}

// Serve starts the server (plain, static-TLS, or autocert) and logs
// a ready message, then waits for graceful shutdown.
func Serve(ctx context.Context, srv *http.Server, certFile, keyFile string) error {
	// ----- log ready (helps during dev & prod) -----
	tlsEnabled := (certFile != "" && keyFile != "") || srv.TLSConfig != nil
	logger.Info("listening", "addr", srv.Addr, "tls", tlsEnabled)

	// ----- run the server in background ------------
	go func() {
		var err error
		switch {
		case certFile != "" && keyFile != "":
			err = srv.ListenAndServeTLS(certFile, keyFile)
		case srv.TLSConfig != nil: // autocert TLSConfig already set in New()
			err = srv.ListenAndServeTLS("", "") // key/cert come from TLSConfig
		default:
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
		}
	}()

	// ----- block until ctx cancelled or SIGINT/SIGTERM -----
	return Graceful(ctx, srv)
}

/* ---------------------------------------------------------------------- */

func pickPort(cfgPort, fallback int) string {
	if cfgPort == 0 {
		return fmt.Sprintf("%d", fallback)
	}
	return fmt.Sprintf("%d", cfgPort)
}
