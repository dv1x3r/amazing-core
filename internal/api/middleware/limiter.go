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
	visitors    map[string]*visitor
	mutex       *sync.RWMutex
	rate        rate.Limit
	burst       int
	expiresIn   time.Duration
	lastCleanup time.Time
}

func RateLimiter(rate rate.Limit, burst int, expiresIn time.Duration) Middleware {
	l := &rateLimiter{
		visitors:  map[string]*visitor{},
		mutex:     &sync.RWMutex{},
		rate:      rate,
		burst:     burst,
		expiresIn: expiresIn,
	}
	return l.middleware
}

func (l *rateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()

	// lazy cleanup
	if now.Sub(l.lastCleanup) > l.expiresIn {
		for k, v := range l.visitors {
			if now.Sub(v.lastSeen) > l.expiresIn {
				delete(l.visitors, k)
			}
		}
		l.lastCleanup = now
	}

	if v, exists := l.visitors[ip]; exists {
		v.lastSeen = now
		return v.limiter
	}

	limiter := rate.NewLimiter(l.rate, l.burst)
	l.visitors[ip] = &visitor{
		limiter:  limiter,
		lastSeen: now,
	}

	return limiter
}

func (l *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Context().Value(IPExtractorKey).(string)
		if limiter := l.getLimiter(ip); !limiter.Allow() {
			res := w2.NewErrorResponse(http.StatusText(http.StatusTooManyRequests))
			res.Write(w, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
