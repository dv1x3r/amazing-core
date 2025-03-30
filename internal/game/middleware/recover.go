package middleware

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
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
