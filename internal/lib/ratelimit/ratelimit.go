package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Limiter is a rate limiter per key, e.g. per IP.
type Limiter struct {
	mu          sync.Mutex
	visitors    map[string]*visitor
	rate        rate.Limit
	burst       int
	expiresIn   time.Duration
	lastCleanup time.Time
}

// New creates a limiter, a key unseen for expiresIn starts with a full burst again.
func New(r rate.Limit, burst int, expiresIn time.Duration) *Limiter {
	return &Limiter{
		visitors:  map[string]*visitor{},
		rate:      r,
		burst:     burst,
		expiresIn: expiresIn,
	}
}

// Allow reports whether the key may make one more request now.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

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

	v, ok := l.visitors[key]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.visitors[key] = v
	}
	v.lastSeen = now
	return v.limiter.Allow()
}
