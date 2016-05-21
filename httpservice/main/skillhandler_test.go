package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThatTheSkillHandlerReturnsJSON(t *testing.T) {
	mda := &mockDataAccess{
		listSkillTagsResponse: func() ([]string, error) {
			return []string{"a", "b"}, nil
		},
	}

	sessionFactory := func(w http.ResponseWriter, r *http.Request) Session {
		return &mockSession{
			validateSessionValidResponse:        true,
			validateSessionEmailAddressResponse: "a-h@github.com",
			validateSessionWasCalled:            false,
		}
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/skill/", nil)

	sh := NewSkillHandler(mda, sessionFactory)
	sh.ServeHTTP(w, r)

	if !strings.Contains(w.HeaderMap["Content-Type"][0], "application/json") {
		t.Fatal("The skill handler should return JSON.")
	}

	expected := `["a","b"]`
	actual := strings.TrimSpace(w.Body.String())

	if actual != expected {
		t.Errorf("Expected JSON to be %s, was %s", expected, actual)
	}
}
