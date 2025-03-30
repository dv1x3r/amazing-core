package middleware

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/w2go/w2"
)

func Protected(service *auth.Service) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := service.GetSessionUsername(w, r); !ok {
				w2.NewErrorResponse(http.StatusText(http.StatusUnauthorized)).Write(w, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
