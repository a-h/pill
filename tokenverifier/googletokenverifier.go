package tokenverifier

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// A GoogleTokenVerifier verifies Tokens with Google.
type GoogleTokenVerifier struct{}

// ValidateToken retrieves a claim from Google and validates it using Google's rules.
func (verifier GoogleTokenVerifier) ValidateToken(idToken string) (claim *Claim, err error) {
	claim, err = verifier.GetClaim(idToken)

	if err != nil {
		return claim, err
	}

	_, err = verifier.IsClaimValid(claim)

	return claim, err
}

// GetClaim returns a Claim from Google, using the id_token presented by the
// Google authentication system.
func (verifier GoogleTokenVerifier) GetClaim(idToken string) (*Claim, error) {
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + url.QueryEscape(idToken)
	body, err := getResponse(url)

	if err != nil {
		return nil, err
	}

	claim, err := NewClaim(body)

	if err != nil {
		log.Print("Unable to understand the claim received from Google.", err)
		return nil, err
	}

	return claim, err
}

func getResponse(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print("Failed to retrieve the URL from google.", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("Failed to receive the claim from Google.", err)
		return nil, err
	}

	return body, nil
}

// IsClaimValid validates a claim by checking that it's not expired, the issuer
// was Google and that the user's email address has been verified.
func (verifier GoogleTokenVerifier) IsClaimValid(claim *Claim) (ok bool, err error) {
	expiry, expiryErr := strconv.Atoi(claim.Expiry)
	emailVerified, emailVerifiedErr := strconv.ParseBool(claim.EmailVerified)

	validation := []struct {
		name               string
		validationFunction func() bool
	}{
		{"email ok", func() bool { return claim.Email != "" }},
		{"email verified ok", func() bool { return emailVerifiedErr == nil && emailVerified }},
		{"expiry is number", func() bool { return expiryErr == nil }},
		{"expiry ok", func() bool { return time.Unix(int64(expiry), 0).After(time.Now()) }},
		{"issuer ok", func() bool {
			return claim.Issuer == "https://accounts.google.com" || claim.Issuer == "accounts.google.com"
		}},
	}

	var errorMessage bytes.Buffer
	ok = true
	for _, v := range validation {
		valid := v.validationFunction()
		errorMessage.WriteString(v.name + " " + strconv.FormatBool(valid) + "\n")
		ok = ok && valid
	}

	if ok {
		return ok, nil
	}

	return ok, errors.New(errorMessage.String())
}
