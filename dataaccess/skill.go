package dataaccess

// Skill stores information about skills.
type Skill struct {
	Skill string `json:"skill"`
	// Level represents the answer to the question "What is your level of expertise?"
	Level DreyfusLevel `json:"level"`
	// Interest represents the answer to the question "You are interested in using this skill for work."
	Interest LikertScale `json:"interest"`
}
