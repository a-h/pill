package dataaccess

import "time"

// SkillLevel holds a number of skills.
type SkillLevel struct {
	Date   time.Time `json:"date"`
	Skills []Skill   `json:"skills"`
}
