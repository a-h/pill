package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/pill/tokenverifier"
)

func TestThatAValidSessionRedirectsToTheProfile(t *testing.T) {
	ms := &mockSession{
		validateSessionValidResponse:        true,
		validateSessionEmailAddressResponse: "a-h@github.com",
		startSessionWasCalled:               false}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	vf := tokenverifier.NewTestTokenVerifier(&tokenverifier.Claim{}, nil)

	lh := NewLoginHandler(sf, vf)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/login", nil)

	lh.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Error("The HTTP status code returned should be StatusFound because the user was redirected.")
	}

	if w.HeaderMap.Get("Location") != "/profile/" {
		t.Error("The user should have been redirected to the profile URL.")
	}

	if ms.startSessionWasCalled {
		t.Error("The session should not have been started, because it's already valid.")
	}
}

func TestThatAnInvalidSessionRendersTheLoginView(t *testing.T) {
	ms := &mockSession{
		validateSessionValidResponse:        false,
		validateSessionEmailAddressResponse: "",
		startSessionWasCalled:               false}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	vf := tokenverifier.NewTestTokenVerifier(&tokenverifier.Claim{}, nil)

	lh := NewLoginHandler(sf, vf)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/login", nil)

	lh.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("The HTTP status code returned should be StatusOK because we're accessing the logon screen.")
	}

	body := w.Body.String()

	if !strings.Contains(body, "Login with your Google Account") {
		t.Error("The login view was not rendered.")
	}
}
