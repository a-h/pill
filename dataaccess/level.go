package dataaccess

// The DreyfusLevel of skill, from 1 to 5.
type DreyfusLevel int

const (
	// Novice people might have read about it.
	Novice = 1
	// Competent people might have used it.
	Competent = 2
	// Proficient people might have used it a lot.
	Proficient = 3
	// Expert people might have written about it.
	Expert = 4
	// Master level people might have taught others about it.
	Master = 5
)
