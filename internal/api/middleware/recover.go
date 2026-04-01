package middleware

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := wrap.Panic(func() error {
				next.ServeHTTP(w, r)
				return nil
			})

			if err != nil {
				if lw, ok := w.(interface{ SetError(error) }); ok {
					lw.SetError(err)
				}
				res := w2.NewErrorResponse(http.StatusText(http.StatusInternalServerError))
				res.Write(w, http.StatusInternalServerError)
			}
		})
	}
}
