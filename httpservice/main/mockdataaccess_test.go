package main

import (
	"testing"

	"github.com/a-h/pill/dataaccess"
)

type mockDataAccess struct {
	getProfileResponse                func(string) (*dataaccess.Profile, bool, error)
	getProfileCallCount               int
	updateProfileResponse             func(update *dataaccess.ProfileUpdate) (*dataaccess.Profile, error)
	updateProfileCallCount            int
	listSkillTagsResponse             func() ([]string, error)
	listSkillTagsCallCount            int
	addSkillTagsResponse              func(tags []string) error
	addSkillTagsCallCount             int
	deleteProfileResponse             func(emailAddress string) (bool, error)
	deleteProfileCallCount            int
	listProfilesResponse              func() ([]dataaccess.Profile, error)
	listProfilesCallCount             int
	deleteSkillTagsResponse           func(tags []string) error
	deleteSkillTagsCallCount          int
	getOrCreateConfigurationResponse  func() (dataaccess.Configuration, error)
	getOrCreateConfigurationCallCount int
	deleteConfigurationResponse       func() error
	deleteConfigurationCallCount      int
}

func (da *mockDataAccess) GetProfile(emailAddress string) (*dataaccess.Profile, bool, error) {
	da.getProfileCallCount++
	return da.getProfileResponse(emailAddress)
}

func (da *mockDataAccess) UpdateProfile(update *dataaccess.ProfileUpdate) (*dataaccess.Profile, error) {
	da.updateProfileCallCount++
	return da.updateProfileResponse(update)
}

func (da *mockDataAccess) ListSkillTags() ([]string, error) {
	da.listSkillTagsCallCount++
	return da.listSkillTagsResponse()
}

func (da *mockDataAccess) AddSkillTags(tags []string) error {
	da.addSkillTagsCallCount++
	return da.addSkillTagsResponse(tags)
}

func (da *mockDataAccess) DeleteProfile(emailAddress string) (bool, error) {
	da.deleteProfileCallCount++
	return da.deleteProfileResponse(emailAddress)
}

func (da *mockDataAccess) ListProfiles(domain string) ([]dataaccess.Profile, error) {
	da.listProfilesCallCount++
	return da.listProfilesResponse()
}

func (da *mockDataAccess) DeleteSkillTags(tags []string) error {
	da.deleteSkillTagsCallCount++
	return da.deleteSkillTagsResponse(tags)
}

func (da *mockDataAccess) GetOrCreateConfiguration() (dataaccess.Configuration, error) {
	da.getOrCreateConfigurationCallCount++
	return da.getOrCreateConfigurationResponse()
}

func (da *mockDataAccess) DeleteConfiguration() error {
	da.deleteConfigurationCallCount++
	return da.deleteConfigurationResponse()
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
