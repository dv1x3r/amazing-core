package webauth

import (
	"crypto/subtle"
	"net/http"

	"github.com/gorilla/sessions"
)

type Authenticator struct {
	session  sessions.Store
	username string
	password string
}

func NewAuthenticator(secure bool, sessionKey, username, password string) *Authenticator {
	session := sessions.NewCookieStore([]byte(sessionKey))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 14,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &Authenticator{
		session:  session,
		username: username,
		password: password,
	}
}

func (a *Authenticator) AuthenticateSession(w http.ResponseWriter, r *http.Request, username, password string) (bool, error) {
	validUsername := subtle.ConstantTimeCompare([]byte(username), []byte(a.username)) == 1
	validPassword := subtle.ConstantTimeCompare([]byte(password), []byte(a.password)) == 1
	if !validUsername || !validPassword {
		return false, nil
	}
	sess, _ := a.session.Get(r, "session")
	sess.Values["username"] = username
	err := sess.Save(r, w)
	return err == nil, err
}

func (a *Authenticator) DeauthenticateSession(w http.ResponseWriter, r *http.Request) error {
	sess, _ := a.session.Get(r, "session")
	delete(sess.Values, "username")
	return sess.Save(r, w)
}

func (a *Authenticator) GetSessionUsername(w http.ResponseWriter, r *http.Request) (string, bool) {
	sess, _ := a.session.Get(r, "session")
	username, ok := sess.Values["username"].(string)
	return username, ok
}

func (a *Authenticator) RefreshSession(w http.ResponseWriter, r *http.Request) error {
	sess, _ := a.session.Get(r, "session")
	return sess.Save(r, w)
}
