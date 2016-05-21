package main

import (
	"log"
	"net/http"

	"github.com/a-h/pill/tokenverifier"
)

// The LoginHandler renders the logon screen and redirects to the profile, once
// the session has been setup.
type LoginHandler struct {
	getSession    func(w http.ResponseWriter, r *http.Request) Session
	TokenVerifier tokenverifier.TokenVerifier
}

// NewLoginHandler creates an instance of the LoginHandler.
func NewLoginHandler(sessionFactory func(w http.ResponseWriter, r *http.Request) Session, verifier tokenverifier.TokenVerifier) *LoginHandler {
	return &LoginHandler{sessionFactory, verifier}
}

func (handler LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleLoginGet(w, r, handler)
	} else {
		handleLoginPost(w, r, handler)
	}
}

func handleLoginGet(w http.ResponseWriter, r *http.Request, handler LoginHandler) {
	log.Print("Handling login get.")

	valid, _ := handler.getSession(w, r).ValidateSession()
	if valid {
		log.Print("The session is valid, redirecting to /profile/")
		http.Redirect(w, r, "/profile/", http.StatusFound)
		return
	}

	log.Print("Rendering the login template.")
	renderTemplate(w, "login.html", nil)
}

func handleLoginPost(w http.ResponseWriter, r *http.Request, handler LoginHandler) {
	r.ParseForm()
	idToken := r.FormValue("id_token")

	claim, err := handler.TokenVerifier.ValidateToken(idToken)

	if err != nil {
		log.Printf("The claim %s is invalid. With error message %s", idToken, err.Error())
		http.Error(w, "The presented claim is invalid.", http.StatusInternalServerError)
		return
	}

	handler.getSession(w, r).StartSession(claim.Email)

	http.Redirect(w, r, "/profile/", http.StatusFound)
}
