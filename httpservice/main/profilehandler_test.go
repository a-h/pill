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
		startSessionWasCalled:               false,
		validateSessionWasCalled:            false,
	}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	ph := NewProfileHandler(nil, sf)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/profile", nil)

	ph.ServeHTTP(w, r)

	if !ms.validateSessionWasCalled {
		t.Error("The session was not validated.")
	}
}
