// middleware/limit.go
package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/*───────────────────────────────────────────────────────────────────────────────
  Body-size limit
───────────────────────────────────────────────────────────────────────────────*/

// MaxBytes wraps h with http.MaxBytesReader (per request).
func MaxBytes(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

/*───────────────────────────────────────────────────────────────────────────────
  IP rate-limit  (token bucket per remote IP)
───────────────────────────────────────────────────────────────────────────────*/

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

		// expire after TTL
		time.AfterFunc(l.ttl, func() {
			l.mu.Lock()
			delete(l.pool, ip)
			l.mu.Unlock()
		})
	}
	return lim
}

// RateLimit returns middleware that allows `events` per `per` with `burst` size
// and sets standard headers:
//
//   - X-RateLimit-Limit       (events)
//   - X-RateLimit-Burst       (burst)
//   - X-RateLimit-Remaining   (tokens left before this request)
//   - Retry-After             (seconds) on 429 responses
func RateLimit(events int, per time.Duration, burst int) func(http.Handler) http.Handler {
	l := newIPLimiter(rate.Every(per/time.Duration(events)), burst, 10*time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			lim := l.get(ip)

			// Headers visible to every client response
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(events))
			w.Header().Set("X-RateLimit-Burst", strconv.Itoa(burst))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(lim.Tokens())))

			if !lim.Allow() {
				retry := lim.Reserve().Delay()
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retry.Seconds()))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
