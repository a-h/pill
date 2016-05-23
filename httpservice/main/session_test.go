package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestThatInvalidSessionsRedirectToTheHomePage(t *testing.T) {
	w := httptest.NewRecorder()
	redirectURL, _ := url.Parse("http://example.com/login")
	r, _ := http.NewRequest("GET", "http://example.com/secret_area", nil)
	s := NewGorillaSession(w, r, []byte("random_data"), false, *redirectURL)

	result, _ := s.ValidateSession()

	if result == true {
		t.Error("The session is not valid, because no cookie exists.")
	}

	if w.Code != http.StatusFound {
		t.Error("The HTTP status code returned should be StatusFound because the user was redirected.")
	}

	if w.HeaderMap.Get("Location") != redirectURL.String() {
		t.Error("The user should have been redirected to the login URL.")
	}
}

func TestThatLoginsDoNotLoop(t *testing.T) {
	w := httptest.NewRecorder()
	redirectURL, _ := url.Parse("/")
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	s := NewGorillaSession(w, r, []byte("random_data"), false, *redirectURL)

	result, _ := s.ValidateSession()

	if result == true {
		t.Error("The session is not valid, because no cookie exists.")
	}

	if w.Code != http.StatusOK {
		t.Error("The HTTP status code returned should be StatusOK. Users shouldn't be redirected to the login screen if that's where they already going.")
	}
}

func TestThatValidSessionsPassThrough(t *testing.T) {
	w := httptest.NewRecorder()
	redirectURL, _ := url.Parse("http://example.com/")
	r, _ := http.NewRequest("GET", "http://example.com/secret_area", nil)

	s := NewGorillaSession(w, r, []byte("random_data"), false, *redirectURL)
	s.StartSession("a-h@github.com")

	result, emailAddress := s.ValidateSession()

	if result == false {
		t.Error("The session is valid, because a cookie exists.")
	}

	if w.Code != http.StatusOK {
		t.Error("The HTTP status code returned should be StatusOK because no action was taken.", w.Code)
	}

	if emailAddress != "a-h@github.com" {
		t.Error("The session should store the user's email address.")
	}
}
