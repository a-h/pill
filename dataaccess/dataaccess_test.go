package dataaccess

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2"
)

func TestThatItIsPossibleToSaveAndUpdateAProfile(t *testing.T) {
	testEmailAddress := "a-h@github.com"
	da := NewMongoDataAccess("mongodb://localhost:27017", "pilltest")
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

func TestSkillTags(t *testing.T) {
	da := NewMongoDataAccess("mongodb://localhost:27017", "pilltest")

	skillTags := []string{"test_tag_" + strconv.Itoa(rand.Int()),
		"test_tag_" + strconv.Itoa(rand.Int())}

	err := da.AddSkillTags(skillTags)

	if err != nil {
		t.Error("Failed to add skill tags. ", err)
	}

	allSkillTags, err := da.ListSkillTags()

	if err != nil {
		t.Error("Failed to retrieve all skill tags. ", err)
	}

	if !containsAll(allSkillTags, skillTags) {
		t.Error("The full list of all skill tags didn't contain the new skill tags.")
	}

	err = da.DeleteSkillTags(skillTags)

	if err != nil {
		t.Error("Failed to delete test skill tags.", err)
	}

	allSkillTags, err = da.ListSkillTags()

	if err != nil {
		t.Error("Failed to retrieve all skill tags (#2). ", err)
	}

	if containsAny(allSkillTags, skillTags) {
		t.Error("After deletion, the test skill tags should not be present in the DB.")
	}
}

func TestContainsAllFunction(t *testing.T) {
	tests := []struct {
		input          []string
		mustContainAll []string
		expected       bool
	}{
		{[]string{"1", "2", "3"}, []string{"1", "2", "3"}, true},
		{[]string{"1", "2"}, []string{"1", "2", "3"}, false},
		{[]string{}, []string{"1"}, false},
		{[]string{"1"}, []string{"1"}, true},
	}

	for _, test := range tests {
		actual := containsAll(test.input, test.mustContainAll)

		if actual != test.expected {
			t.Errorf("containsAll for source %s and mustContainAll %s should have returned %t, but returned %t.",
				test.input, test.mustContainAll, test.expected, actual)
		}
	}
}

func makeMap(slice []string) map[string]bool {
	sliceMap := make(map[string]bool)

	for _, v := range slice {
		sliceMap[v] = true
	}

	return sliceMap
}

func containsAll(source []string, mustContainAllOf []string) bool {
	mustContainAllOfMap := makeMap(mustContainAllOf)

	contains := 0
	for _, a := range source {
		if _, ok := mustContainAllOfMap[a]; ok {
			contains++
		}
	}

	return contains == len(mustContainAllOf)
}

func TestContainsAnyFunction(t *testing.T) {
	tests := []struct {
		input         []string
		containsAnyOf []string
		expected      bool
	}{
		{[]string{"a", "b", "c"}, []string{"a"}, true},
		{[]string{"a", "b", "c"}, []string{"b"}, true},
		{[]string{"a", "b", "c"}, []string{"c"}, true},
		{[]string{"a", "b", "c"}, []string{"d"}, false},
	}

	for _, test := range tests {
		actual := containsAny(test.input, test.containsAnyOf)

		if actual != test.expected {
			t.Errorf("Input %s, parameter %s, expected %t, actual %t", test.input, test.containsAnyOf, test.expected, actual)
		}
	}
}

func containsAny(source []string, containsAnyOf []string) bool {
	mustContainAnyOfMap := makeMap(containsAnyOf)

	for _, a := range source {
		if _, ok := mustContainAnyOfMap[a]; ok {
			return true
		}
	}

	return false
}

func TestThatTagsCanBeCleaned(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"UPPER", "upper"},
		{"space in the Tags", "space-in-the-tags"},
	}

	for _, c := range cases {
		actual := CleanTag(c.in)

		if actual != c.expected {
			t.Errorf("Input '%s', Expected '%s', Actual '%s'", c.in, c.expected, actual)
		}
	}
}

func TestThatConfigurationCanBeRecreated(t *testing.T) {
	da := NewMongoDataAccess("mongodb://localhost:27017", "pilltest")

	// Clean up before testing.
	err := da.DeleteConfiguration()

	if err != nil {
		t.Error("Failed to clean up the configuration collection (#1).", err)
	}

	c1, err := da.GetOrCreateConfiguration()

	if err != nil {
		t.Fatal("Failed to get or create the configuration (#1).")
	}

	err = da.DeleteConfiguration()

	if err != nil {
		t.Error("Failed to clean up the configuration collection (#2).", err)
	}

	c2, err := da.GetOrCreateConfiguration()

	if err != nil {
		t.Error("Failed to get the configuration (#2).")
	}

	if c1.ID != c2.ID {
		t.Error("The ID value of the configuration entry should always be 'configuration'")
	}

	if reflect.DeepEqual(c1, c2) {
		t.Error("After deleting configuration, attempting to create a new configuration entry should result in a random session key being generated.")
		t.Errorf("Key 1: %v", c1)
		t.Errorf("Key 2: %v", c2)
	}
}

func TestThatDeepEqualComparesArrays(t *testing.T) {
	cases := []struct {
		a                []byte
		b                []byte
		expectedAreEqual bool
	}{
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, true},
		{[]byte{1, 2}, []byte{1, 2, 3}, false},
		{[]byte{1, 2, 3}, []byte{1, 2}, false},
		{nil, []byte{1, 2}, false},
	}

	for _, c := range cases {
		actual := reflect.DeepEqual(c.a, c.b)

		if actual != c.expectedAreEqual {
			t.Errorf("For inputs %v and %v, deep equal was expected to return %t, but returned %t.", c.a, c.b, c.expectedAreEqual, actual)
		}
	}
}
