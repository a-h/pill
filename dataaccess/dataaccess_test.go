package dataaccess

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2"
)

func TestThatItIsPossibleToSaveAndUpdateAProfile(t *testing.T) {
	testEmailAddress := "a-h@github.com"
	da := NewMongoDataAccess("mongodb://localhost:27017")
	update := NewProfileUpdate()
	update.Availability = Red
	update.EmailAddress = testEmailAddress

	_, err := da.DeleteProfile(testEmailAddress)

	if err != mgo.ErrNotFound {
		t.Fatal("Failed to clean up the database.", err)
	}

	r1, err := da.UpdateProfile(update)

	if err != nil {
		t.Error("Failed to update profile.", err)
	}

	if r1.EmailAddress != testEmailAddress {
		t.Errorf("Expected an email address of %s, was %s.", testEmailAddress, r1.EmailAddress)
	}

	if len(r1.Skills) > 0 {
		t.Error("Expected the newly created profile to be empty.")
	}

	r2, found, err := da.GetProfile(testEmailAddress)

	if err != nil || !found {
		t.Error("Failed to retrieve a profile.", err)
	}

	differentProperties := []string{}

	if r1.Availability != r2.Availability {
		differentProperties = append(differentProperties, "Availability")
	}

	if r1.EmailAddress != r2.EmailAddress {
		differentProperties = append(differentProperties, "EmailAddress")
	}

	if r1.LastUpdated != r2.LastUpdated {
		differentProperties = append(differentProperties, fmt.Sprintf("LastUpdated %s - %s", r1.LastUpdated, r2.LastUpdated))
	}

	if len(r1.Skills) != len(r2.Skills) {
		differentProperties = append(differentProperties, "Skills (length)")
	}

	if len(differentProperties) > 0 {
		t.Error("When the newly created profile is returned from a get operation, it should be the same as the newly created profile.", differentProperties)
	}

	_, err = da.DeleteProfile(testEmailAddress)

	if err != nil {
		t.Error("Failed to delete the profile.", err)
	}

	_, found, _ = da.GetProfile(testEmailAddress)

	if found {
		t.Error("After deletion, the profile was still there (and it shouldn't be).")
	}
}
