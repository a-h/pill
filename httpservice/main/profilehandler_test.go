package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/pill/dataaccess"
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

func TestThatTheProfileHandlerRendersTheProfileView(t *testing.T) {
	ms := &mockSession{
		validateSessionValidResponse:        true,
		validateSessionEmailAddressResponse: "a-h@github.com",
		startSessionWasCalled:               false,
		validateSessionWasCalled:            false,
	}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	mda := &mockDataAccess{
		getProfileResponse: func(emailAddress string) (profile *dataaccess.Profile, ok bool, err error) {
			return dataaccess.NewProfile(), true, nil
		},
	}

	ph := NewProfileHandler(mda, sf)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/profile", nil)

	ph.ServeHTTP(w, r)

	if !ms.validateSessionWasCalled {
		t.Error("The session must be validated by the handler.")
	}

	if mda.getProfileCallCount != 1 {
		t.Error("The data access code should have been called.")
	}
}
