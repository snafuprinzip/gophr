package main

import (
	"net/http"
	"net/url"
	"time"
)

// Session contains all Session account information
type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

const (
	// Keep users logged in for 3 days
	sessionLength     = 24 * 3 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLength   = 20
)

// NewSession generates a new Session record and attaches corresponding login cookie
func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)
	session := &Session{
		ID:     GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}

	http.SetCookie(w, &cookie)
	return session
}

// RequestSession retrieves the Session from a http Request cookie or returns nil
// if not found
func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}

	// session does not exist in store
	if session == nil {
		return nil
	}

	// delete session from store if it has expired
	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

// Expired checks the expiry date and returns true if the session timeoout
// has been reached
func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

// RequestUser retrieves the User from a http Request cookie or returns nil
// if not found
func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserID == "" {
		return nil
	}

	user, err := globalUserStore.Find(session.UserID)
	if err != nil {
		panic(err)
	}

	return user
}

// RequireLogin checks if RequestUser returns a valid user or if not set's the
// next entry to the requested url and redirects to the login page
func RequireLogin(w http.ResponseWriter, r *http.Request) {
	// pass if user is found
	if RequestUser(r) != nil {
		return
	}

	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))

	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

// FindOrCreateSession looks for an already existing session for this user or
// create a new session if none is found
func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}
	return session
}
