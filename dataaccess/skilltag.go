package dataaccess

// SkillTag names a skill, e.g. "c#", "java"
type SkillTag struct {
	Name string `bson:"_id" json:"name"`
}
