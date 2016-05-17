package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThatItIsNotPossibleToAccessTheProfileWithAnInvalidSession(t *testing.T) {
	ms := &mockSession{
		validateSessionValidResponse:        false,
		validateSessionEmailAddressResponse: "a-h@github.com",
		startSessionWasCalled:               false}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	ph := NewProfileHandler(nil, sf)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/profile", nil)

	ph.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Error("The HTTP status code returned should be StatusFound because the user was redirected.")
	}

	if w.HeaderMap.Get("Location") != "/" {
		t.Error("The user should have been redirected to the login URL, because they don't have a valid session.")
	}

	if ms.startSessionWasCalled {
		t.Error("The session should not have been started, because it's already valid.")
	}
}
