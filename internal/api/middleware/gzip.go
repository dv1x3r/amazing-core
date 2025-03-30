package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	compressor io.Writer
}

func (w *gzipResponseWriter) WriteHeader(c int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(c)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	w.Header().Del("Content-Length")
	return w.compressor.Write(b)
}

func Gzip() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var compress bool
			for _, encoding := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
				if strings.TrimSpace(encoding) == "gzip" {
					compress = true
					break
				}
			}

			// prevent intermediate caches corruption
			w.Header().Add("Vary", "Accept-Encoding")

			// do not compress if not needed
			if !compress || r.Header.Get("Upgrade") != "" {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			r.Header.Del("Accept-Encoding")

			gz := gzip.NewWriter(w)
			defer gz.Close()

			gw := &gzipResponseWriter{ResponseWriter: w, compressor: gz}
			next.ServeHTTP(gw, r)
		})
	}
}
