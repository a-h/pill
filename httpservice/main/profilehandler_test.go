package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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

func TestThatAfterAProfileUpdateTheUserIsRedirectedToTheReport(t *testing.T) {
	ms := &mockSession{
		validateSessionValidResponse:        true,
		validateSessionEmailAddressResponse: "a-h@github.com",
		startSessionWasCalled:               false,
		validateSessionWasCalled:            false,
	}

	sf := func(w http.ResponseWriter, r *http.Request) Session {
		return ms
	}

	var receivedEmailAddress string
	var receivedAvailability dataaccess.RagStatus
	var receivedSkills []dataaccess.Skill

	mda := &mockDataAccess{
		updateProfileResponse: func(update *dataaccess.ProfileUpdate) (*dataaccess.Profile, error) {
			receivedAvailability = update.Availability
			receivedEmailAddress = update.EmailAddress
			receivedSkills = update.Skills

			return dataaccess.NewProfile(), nil
		},
	}

	ph := NewProfileHandler(mda, sf)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("availability", strconv.Itoa(dataaccess.Amber))

	// Add some skills in.
	form.Add("name_1", "C# Development") // This is uppercase on purpose, the handler should lowercase it.
	form.Add("level_1", "5")
	form.Add("interest_1", "3")

	form.Add("name_2", "golang")
	form.Add("level_2", "2")
	form.Add("interest_2", "5")

	r, _ := http.NewRequest("POST", "http://example.com/profile", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ph.ServeHTTP(w, r)

	if !ms.validateSessionWasCalled {
		t.Error("The session must be validated by the handler.")
	}

	if mda.updateProfileCallCount != 1 {
		t.Error("Posting a valid update form to the handler should result in a profile update.")
	}

	if receivedEmailAddress != "a-h@github.com" {
		t.Errorf("Expected that the update should apply to a-h@github.com, but was %s.", receivedEmailAddress)
	}

	if receivedAvailability != dataaccess.Amber {
		t.Errorf("Expected that the updated availability was Amber, but %d was received.", receivedAvailability)
	}

	if len(receivedSkills) != 2 {
		t.Errorf("Expected to receive 2 skills in the data update, but only received %d", len(receivedSkills))
	}

	expectedSkill := dataaccess.Skill{
		Skill:    "C# Development",
		Level:    dataaccess.MasterLevel,
		Interest: dataaccess.NeitherAgreeNorDisagree,
	}

	if containsAll(receivedSkills, expectedSkill) {
		t.Error("A skill with 'C#' (uppercase) should not have been found, the handler should have lowercased it and replaced the space with a hyphen.")
	}

	expectedSkill1 := dataaccess.Skill{
		Skill:    "c#-development",
		Level:    dataaccess.MasterLevel,
		Interest: dataaccess.NeitherAgreeNorDisagree,
	}

	expectedSkill2 := dataaccess.Skill{
		Skill:    "golang",
		Level:    dataaccess.CompetentLevel,
		Interest: dataaccess.StronglyAgree,
	}

	if !containsAll(receivedSkills, expectedSkill1, expectedSkill2) {
		t.Errorf("The full list of skills was not retrieved directly from the form, which only received %v", receivedSkills)
	}
}

func containsAll(receivedSkills []dataaccess.Skill, expectedSkills ...dataaccess.Skill) bool {
	found := 0

	for _, skill := range receivedSkills {
		for _, expectedSkill := range expectedSkills {
			if reflect.DeepEqual(skill, expectedSkill) {
				found++
				continue
			}
		}
	}

	return found == len(expectedSkills)
}
