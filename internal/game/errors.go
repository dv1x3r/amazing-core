package game

import (
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/appcode"
)

var (
	// errNoSession is the error response for a connection without a session,
	// e.g. a request the client re-sends after a reconnect before RELOGIN;
	// a plain error would close the connection before RELOGIN.
	errNoSession = wrap.WithGSFError(errors.New("no session for the connection"), appcode.INTERNAL_ERROR)
)
