package middleware

import (
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

func Recover() gsf.Middleware {
	return func(next gsf.HandlerFunc) gsf.HandlerFunc {
		return func(w gsf.ResponseWriter, r *gsf.Request) error {
			return wrap.Panic(func() error {
				return next(w, r)
			})
		}
	}
}
