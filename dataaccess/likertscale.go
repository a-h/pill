package dataaccess

// LikertScale is a scale represented by integer values from 1 to 5.
type LikertScale int

const (
	// StronglyDisagree on the Likert scale
	StronglyDisagree = 1
	// Disagree on the Likert scale
	Disagree = 2
	// NeitherAgreeNorDisagree on the Likert scale
	NeitherAgreeNorDisagree = 3
	// Agree on the Likert scale
	Agree = 4
	// StronglyAgree on the Likert scale
	StronglyAgree = 5
)
