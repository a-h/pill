package dataaccess

import "time"

// ProfileUpdate is used to update a profile.
type ProfileUpdate struct {
	EmailAddress string    `json:"emailAddress"`
	Skills       []Skill   `json:"skills"`
	Availability RagStatus `json:"availability"`
}

// NewProfileUpdate creates an empty profile update.
func NewProfileUpdate() *ProfileUpdate {
	return &ProfileUpdate{}
}

// Profile returns the profile of a person.
type Profile struct {
	EmailAddress  string       `bson:"_id" json:"emailAddress"`
	Skills        []Skill      `json:"skills"`
	Availability  RagStatus    `json:"availability"`
	SkillsHistory []SkillLevel `json:"skillsHistory"`
	Version       int          `json:"version"`
	LastUpdated   time.Time    `json:"lastUpdated"`
	Domain        string       `json:"domain"`
}

// NewProfile creates an empty profile.
func NewProfile() *Profile {
	return &Profile{
		LastUpdated: time.Now(),
	}
}
