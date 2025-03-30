package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/dv1x3r/w2go/w2"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu        *sync.RWMutex
	visitors  map[string]*visitor
	rate      rate.Limit
	burst     int
	expiresIn time.Duration
}

func RateLimiter(rate rate.Limit, burst int, expiresIn time.Duration) Middleware {
	l := &rateLimiter{
		mu:        &sync.RWMutex{},
		visitors:  map[string]*visitor{},
		rate:      rate,
		burst:     burst,
		expiresIn: expiresIn,
	}
	go l.cleanupVisitors()
	return l.middleware
}

func (l *rateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.RLock()
	v, ok := l.visitors[ip]
	l.mu.RUnlock()
	if ok {
		v.lastSeen = time.Now()
		return v.limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	limiter := rate.NewLimiter(l.rate, l.burst)
	l.visitors[ip] = &visitor{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func (l *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.expiresIn {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Context().Value(IPExtractorKey).(string)
		if limiter := l.getLimiter(ip); !limiter.Allow() {
			w2.NewErrorResponse(http.StatusText(http.StatusTooManyRequests)).Write(w, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
