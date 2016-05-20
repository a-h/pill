package main

import (
	"testing"

	"github.com/a-h/pill/dataaccess"
)

type mockDataAccess struct {
	getProfileResponse     func(string) (*dataaccess.Profile, bool, error)
	getProfileCallCount    int
	updateProfileResponse  func(update *dataaccess.ProfileUpdate) (*dataaccess.Profile, error)
	updateProfileCallCount int
	listSkillTagsResponse  func() ([]string, error)
	listSkillTagsCallCount int
	addSkillTagsResponse   func(tags []string) error
	addSkillTagsCallCount  int
	deleteProfileResponse  func(emailAddress string) (bool, error)
	deleteProfileCallCount int
	listProfilesResponse   func() ([]dataaccess.Profile, error)
	listProfilesCallCount  int
}

// GetProfile returns a Profile by the email address of the person.
func (da *mockDataAccess) GetProfile(emailAddress string) (*dataaccess.Profile, bool, error) {
	da.getProfileCallCount++
	return da.getProfileResponse(emailAddress)
}

// UpdateProfile updates a person's profile and returns the newly created
// or updated profile.
func (da *mockDataAccess) UpdateProfile(update *dataaccess.ProfileUpdate) (*dataaccess.Profile, error) {
	da.updateProfileCallCount++
	return da.updateProfileResponse(update)
}

// ListSkillTags lists the skills used before.
func (da *mockDataAccess) ListSkillTags() ([]string, error) {
	da.listSkillTagsCallCount++
	return da.listSkillTagsResponse()
}

// AddSkillTags adds a skill tag to the list.
func (da *mockDataAccess) AddSkillTags(tags []string) error {
	da.addSkillTagsCallCount++
	return da.addSkillTagsResponse(tags)
}

// DeleteProfile removes a profile specified by email address.
func (da *mockDataAccess) DeleteProfile(emailAddress string) (bool, error) {
	da.deleteProfileCallCount++
	return da.deleteProfileResponse(emailAddress)
}

// ListProfiles lists all of the profiles stored in the database.
func (da *mockDataAccess) ListProfiles() ([]dataaccess.Profile, error) {
	da.listProfilesCallCount++
	return da.listProfilesResponse()
}

func TestMockRecordsIncrementsAndExecutesFunctions(t *testing.T) {
	mda := &mockDataAccess{
		getProfileResponse: func(string) (*dataaccess.Profile, bool, error) { return nil, false, nil },
	}

	mda.GetProfile("test")

	if mda.getProfileCallCount != 1 {
		t.Error("The profile call count was not incremented.")
	}
}
