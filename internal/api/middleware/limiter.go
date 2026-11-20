package middleware

import (
	"net/http"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/ratelimit"
	"github.com/dv1x3r/w2go/w2"
	"golang.org/x/time/rate"
)

func RateLimiter(limit rate.Limit, burst int, expiresIn time.Duration) Middleware {
	limiter := ratelimit.New(limit, burst, expiresIn)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Context().Value(IPExtractorKey).(string)
			if !limiter.Allow(ip) {
				res := w2.NewErrorResponse(http.StatusText(http.StatusTooManyRequests))
				res.Write(w, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
