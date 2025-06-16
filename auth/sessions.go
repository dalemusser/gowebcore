package auth

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

type Session struct {
	sc *securecookie.SecureCookie
}

func NewSession(hashKey, blockKey []byte) *Session {
	return &Session{sc: securecookie.New(hashKey, blockKey)}
}

func (s *Session) Set(w http.ResponseWriter, name string, value map[string]any) error {
	encoded, err := s.sc.Encode(name, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func (s *Session) Get(r *http.Request, name string, dst *map[string]any) error {
	c, err := r.Cookie(name)
	if err != nil {
		return err
	}
	return s.sc.Decode(name, c.Value, dst)
}

func (s *Session) Clear(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
