package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Service struct {
	session  sessions.Store
	username string
	password string
}

func NewService(secure bool, sessionKey string, username string, password string) *Service {
	session := sessions.NewCookieStore([]byte(sessionKey))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 14,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &Service{
		session:  session,
		username: username,
		password: password,
	}
}
