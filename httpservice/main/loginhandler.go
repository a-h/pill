package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// The LoginHandler renders the logon screen and redirects to the profile, once
// the session has been setup.
type LoginHandler struct {
	getSession func(w http.ResponseWriter, r *http.Request) Session
}

// NewLoginHandler creates an instance of the LoginHandler.
func NewLoginHandler(sessionFactory func(w http.ResponseWriter, r *http.Request) Session) *LoginHandler {
	return &LoginHandler{sessionFactory}
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

	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + url.QueryEscape(idToken)
	body, err := getResponse(url)

	if err != nil {
		http.Error(w, "Failed to receive a response from the Google validation service.", http.StatusInternalServerError)
		return
	}

	claim := &claim{}
	err = json.Unmarshal(body, claim)

	if err != nil {
		log.Print("Unable to understand the claim received from Google.", err)
		http.Error(w, "Unable to connect to Google to validate the logon token.", http.StatusInternalServerError)
		return
	}

	ok, errorMessage := isClaimValid(claim)

	if !ok {
		log.Print("The claim is invalid. ", errorMessage)
		http.Error(w, "The presened claim is invalid.", http.StatusInternalServerError)
		return
	}

	handler.getSession(w, r).StartSession(claim.Email)

	http.Redirect(w, r, "/profile/", http.StatusFound)
}

func isClaimValid(claim *claim) (ok bool, msg string) {
	expiry, expiryErr := strconv.Atoi(claim.Expiry)
	emailVerified, emailVerifiedErr := strconv.ParseBool(claim.EmailVerified)

	validation := map[string]bool{
		"email ok":          claim.Email != "",
		"email verified ok": emailVerifiedErr == nil && emailVerified,
		"expiry is number":  expiryErr == nil,
		"expiry ok":         time.Unix(int64(expiry), 0).After(time.Now()),
		"issuer ok":         claim.Issuer == "https://accounts.google.com" || claim.Issuer == "accounts.google.com",
	}

	var errorMessage bytes.Buffer
	ok = true
	for k, v := range validation {
		errorMessage.WriteString(string(k) + " " + strconv.FormatBool(v) + "\n")
		ok = ok && v
	}

	return ok, errorMessage.String()
}

func getResponse(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print("Failed to retrieve the URL from google.", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("Failed to receive the claim from Google.", err)
		return nil, err
	}

	return body, nil
}

type claim struct {
	// The issuer, should be "https://accounts.google.com" or "accounts.google.com"
	Issuer string `json:"iss"`
	// The expiry, e.g. "1433981953". Should not be in the past.
	Expiry        string `json:"exp"`
	Email         string `json:"email"` // e.g. "testuser@gmail.com",
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`        // e.g. "Test User",
	Picture       string `json:"picture"`     // e.g. "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
	GivenName     string `json:"given_name"`  // e.g. "Test"
	FamilyName    string `json:"family_name"` // e.g. "User"
}
