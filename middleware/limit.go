// middleware/limit.go
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// -----------------------------------------------------------------------------
// Body-size limit
// -----------------------------------------------------------------------------

// MaxBytes wraps h with http.MaxBytesReader (per request).
func MaxBytes(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// -----------------------------------------------------------------------------
// Simple IP rate-limit  (token bucket per remote IP)
// -----------------------------------------------------------------------------

type ipLimiter struct {
	mu   sync.Mutex
	pool map[string]*rate.Limiter
	r    rate.Limit
	b    int
	ttl  time.Duration
}

func newIPLimiter(r rate.Limit, b int, ttl time.Duration) *ipLimiter {
	return &ipLimiter{
		pool: make(map[string]*rate.Limiter),
		r:    r, b: b, ttl: ttl,
	}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	lim, ok := l.pool[ip]
	if !ok {
		lim = rate.NewLimiter(l.r, l.b)
		l.pool[ip] = lim

		// Expire after TTL to avoid unbounded map.
		time.AfterFunc(l.ttl, func() {
			l.mu.Lock()
			delete(l.pool, ip)
			l.mu.Unlock()
		})
	}
	return lim
}

// RateLimit returns a middleware that allows burst `burst`
// and average `events` per `per` duration (e.g. 10 req / sec).
func RateLimit(events int, per time.Duration, burst int) func(http.Handler) http.Handler {
	l := newIPLimiter(rate.Every(per/time.Duration(events)), burst, 10*time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !l.get(ip).Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
