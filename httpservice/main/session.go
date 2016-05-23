package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
)

// Session determines how a user is logged in to the system.
type Session interface {
	// ValidateSession validates a session and returns the email address of the
	// user.
	ValidateSession() (isValid bool, emailAddress string)
	StartSession(emailAddress string)
}

// A GorillaSession uses the Gorilla framework to manage the session.
type GorillaSession struct {
	store    sessions.CookieStore
	w        http.ResponseWriter
	r        *http.Request
	loginURL url.URL
}

const sessionName string = "pill-session-cookie"

// NewGorillaSession creates a Session which uses Gorilla.
func NewGorillaSession(w http.ResponseWriter, r *http.Request, encryptionKey []byte, setSecureFlag bool, loginURL url.URL) *GorillaSession {
	store := sessions.NewCookieStore(encryptionKey)
	store.Options = &sessions.Options{
		HttpOnly: true,
		Secure:   setSecureFlag,
	}

	return &GorillaSession{
		store:    *store,
		w:        w,
		r:        r,
		loginURL: loginURL,
	}
}

// StartSession starts off a session by adding the emailAddress value to an
// encrypted cookie.
func (gs GorillaSession) StartSession(emailAdress string) {
	session, err := gs.store.Get(gs.r, sessionName)

	if err != nil {
		http.Error(gs.w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["emailAddress"] = emailAdress
	session.Save(gs.r, gs.w)
}

// ValidateSession checks whether the session is valid. If it isn't, it will
// redirect the user to the logon screen.
func (gs GorillaSession) ValidateSession() (isValid bool, emailAddress string) {
	log.Print("Validating the session.")
	session, err := gs.store.Get(gs.r, sessionName)
	if err != nil {
		print("Failed to get the cookie from the store.")
		http.Error(gs.w, err.Error(), http.StatusInternalServerError)
		return false, ""
	}

	ea, ok := session.Values["emailAddress"].(string)

	if !ok || ea == "" {
		log.Printf("Failed to recover the email address {ea: %s, ok: %t}. Considering redirecting to %s", ea, ok, gs.loginURL.String())

		log.Printf("The incoming URL was %s.", gs.r.URL.Path)

		if gs.r.URL.Path == gs.loginURL.String() {
			log.Print("Not redirecting because the user is at the logon screen.")
		} else {
			http.Redirect(gs.w, gs.r, gs.loginURL.String(), http.StatusFound)
		}
		return false, ea
	}

	log.Printf("The session is valid for user %s", ea)
	return true, ea
}
