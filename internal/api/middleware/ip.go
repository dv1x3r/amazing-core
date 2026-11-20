package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const IPExtractorKey = "ip"

func IPExtractor() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIPDirect(r)

			// Any client can send the client IP headers, only a reverse proxy on the same host
			// or in the private network (e.g. another docker container) is trusted.
			if parsed := net.ParseIP(ip); parsed != nil && (parsed.IsLoopback() || parsed.IsPrivate()) {
				if forwarded := extractIPFromXFFHeader(r); forwarded != "" {
					ip = forwarded
				} else if realIP := extractIPFromRealIPHeader(r); realIP != "" {
					ip = realIP
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), IPExtractorKey, ip))

			next.ServeHTTP(w, r)
		})
	}
}

func extractIPDirect(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func extractIPFromXFFHeader(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	return ""
}

func extractIPFromRealIPHeader(r *http.Request) string {
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip := strings.TrimSpace(xrip)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}
