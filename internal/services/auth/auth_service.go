package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/config"

	"github.com/gorilla/sessions"
)

type Service struct {
	session sessions.Store
}

func NewService(session sessions.Store) *Service {
	return &Service{session: session}
}

func (s *Service) AuthenticateSession(w http.ResponseWriter, r *http.Request, form AdminLoginForm) (bool, error) {
	validUsername := subtle.ConstantTimeCompare([]byte(form.Username), []byte(config.Get().Secure.Auth.Username)) == 1
	validPassword := subtle.ConstantTimeCompare([]byte(form.Password), []byte(config.Get().Secure.Auth.Password)) == 1
	if !validUsername || !validPassword {
		return false, nil
	}

	sess, _ := s.session.Get(r, "session")
	sess.Values["username"] = form.Username
	err := sess.Save(r, w)
	return err == nil, err
}

func (s *Service) DeauthenticateSession(w http.ResponseWriter, r *http.Request) error {
	sess, _ := s.session.Get(r, "session")
	delete(sess.Values, "username")
	return sess.Save(r, w)
}

func (s *Service) GetSessionUsername(w http.ResponseWriter, r *http.Request) (string, bool) {
	sess, _ := s.session.Get(r, "session")
	username, ok := sess.Values["username"].(string)
	return username, ok
}

func (s *Service) RefreshSession(w http.ResponseWriter, r *http.Request) error {
	sess, _ := s.session.Get(r, "session")
	return sess.Save(r, w)
}
