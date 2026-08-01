package httpserver

import (
	"errors"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
)

var errTooManyRequests = errors.New("too many requests")

// ipRateLimiter enforces a per-client-IP rate limit on endpoints that are
// public and cost an attacker nothing to hit (admin login, lead
// submission) — see architecture review finding #8. Keyed by the trusted
// client IP set by middleware.ClientIPFromXFFTrustedProxies in NewRouter,
// never raw r.RemoteAddr (spoofable behind a proxy hop).
type ipRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	return &ipRateLimiter{limiters: make(map[string]*rate.Limiter), r: r, b: b}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	limiter, ok := l.limiters[ip]
	if !ok {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = limiter
	}
	l.mu.Unlock()
	return limiter.Allow()
}

// middleware fails OPEN (no limiting applied) when the trusted client IP
// can't be determined — an XFF-chain anomaly should never lock out admin
// login, and an attacker can't manufacture that condition themselves (the
// one trusted proxy hop is the one adding the header).
func (l *ipRateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := middleware.GetClientIP(r.Context())
		if ip != "" && !l.allow(ip) {
			writeError(w, http.StatusTooManyRequests, errTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
