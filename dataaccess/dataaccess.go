package dataaccess

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

// The DataAccess interface defines how data is written to the data store.
type DataAccess interface {
	UpdateProfile(update *ProfileUpdate) (*Profile, error)
	GetProfile(emailAddress string) (*Profile, bool, error)
	ListSkillTags() ([]string, error)
	AddSkillTags(tags []string) error
	DeleteProfile(emailAddress string) (bool, error)
	ListProfiles() ([]Profile, error)
}

// MongoDataAccess provides access to the data structures.
type MongoDataAccess struct {
	connectionString string
}

// NewMongoDataAccess creates an instance of the MongoDataAccess type.
func NewMongoDataAccess(connectionString string) DataAccess {
	return &MongoDataAccess{connectionString}
}

// GetProfile returns a Profile by the email address of the person.
func (da MongoDataAccess) GetProfile(emailAddress string) (*Profile, bool, error) {
	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB.", err)
		return nil, false, err
	}
	defer session.Close()

	c := session.DB("pill").C("profiles")

	result := NewProfile()
	result.EmailAddress = emailAddress
	err = c.FindId(emailAddress).One(result)

	if err == mgo.ErrNotFound {
		log.Printf("Failed to find a profile with email %s.", emailAddress)
		return result, false, nil
	}

	return result, true, nil
}

// UpdateProfile updates a person's profile and returns the newly created
// or updated profile.
func (da MongoDataAccess) UpdateProfile(update *ProfileUpdate) (*Profile, error) {
	log.Printf("Updating profile for %s", update.EmailAddress)

	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB.", err)
		return nil, err
	}
	defer session.Close()

	c := session.DB("pill").C("profiles")

	profile, found, err := da.GetProfile(update.EmailAddress)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	if found {
		log.Printf("Found existing profile for %s", update.EmailAddress)
	} else {
		log.Printf("New profile found for %s", update.EmailAddress)
	}

	if len(profile.Skills) > 0 {
		// Move current skills to history, if it's an update to an existing profile.
		sl := SkillLevel{
			Date:   profile.LastUpdated,
			Skills: profile.Skills,
		}

		profile.SkillsHistory = append(profile.SkillsHistory, sl)
	}
	profile.Skills = update.Skills
	profile.Availability = update.Availability
	profile.Version++
	profile.LastUpdated = time.Unix(time.Now().Unix(), 0)

	_, err = c.UpsertId(profile.EmailAddress, profile)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return profile, nil
}

// ListSkillTags lists the skills used before.
func (da MongoDataAccess) ListSkillTags() ([]string, error) {
	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB. ", err)
		return nil, err
	}
	defer session.Close()

	c := session.DB("pill").C("skills")

	var results []SkillTag
	err = c.Find(nil).All(&results)

	if err != nil {
		log.Print("Failed to list skill tags. ", err)
		return nil, nil
	}

	skillTags := make([]string, len(results), len(results))
	for idx, tag := range results {
		skillTags[idx] = tag.Name
	}

	return skillTags, nil
}

// AddSkillTags adds a skill tag to the list.
func (da MongoDataAccess) AddSkillTags(tags []string) error {
	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB.", err)
		return err
	}
	defer session.Close()

	c := session.DB("pill").C("skills")

	documents := make([]SkillTag, len(tags), len(tags))
	for idx, tag := range tags {
		documents[idx] = SkillTag{tag}
	}

	return c.Insert(documents)
}

// DeleteProfile removes a profile specified by email address.
func (da MongoDataAccess) DeleteProfile(emailAddress string) (bool, error) {
	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB.", err)
		return false, err
	}
	defer session.Close()

	err = session.DB("pill").C("profiles").RemoveId(emailAddress)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (da MongoDataAccess) ListProfiles() ([]Profile, error) {
	session, err := mgo.Dial(da.connectionString)
	if err != nil {
		log.Print("Failed to connect to MongoDB.", err)
		return nil, err
	}
	defer session.Close()

	var results []Profile
	err = session.DB("pill").C("profiles").Find(nil).All(&results)

	if err != nil {
		log.Print("Failed to list profiles.", err)
		return nil, err
	}

	return results, nil
}