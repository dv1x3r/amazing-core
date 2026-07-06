package middleware

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/webauth"
	"github.com/dv1x3r/w2go/w2"
)

func Protected(authenticator *webauth.Authenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := authenticator.GetSessionUsername(w, r); !ok {
				res := w2.NewErrorResponse(http.StatusText(http.StatusUnauthorized))
				res.Write(w, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
