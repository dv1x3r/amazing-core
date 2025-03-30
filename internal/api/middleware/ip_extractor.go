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
			ip := extractIPFromXFFHeader(r)

			if ip == "" {
				ip = extractIPFromRealIPHeader(r)
			}

			if ip == "" {
				ip = extractIPDirect(r)
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
