package dataaccess

// The DreyfusLevel of skill, from 1 to 5.
type DreyfusLevel int

const (
	// NoviceLevel people might have read about it.
	NoviceLevel = 1
	// CompetentLevel people might have used it.
	CompetentLevel = 2
	// ProficientLevel people might have used it a lot.
	ProficientLevel = 3
	// ExpertLevel people might have written about it.
	ExpertLevel = 4
	// MasterLevel people might have taught others about it.
	MasterLevel = 5
)
