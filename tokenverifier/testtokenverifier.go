package tokenverifier

// TestTokenVerifier is a test implementation of the TokenVerifier.
type TestTokenVerifier struct {
	claimToReturn *Claim
	errorToReturn error
}

// NewTestTokenVerifier creates a new test token.
func NewTestTokenVerifier(claimToReturn *Claim, errorToReturn error) *TestTokenVerifier {
	return &TestTokenVerifier{claimToReturn, errorToReturn}
}

// ValidateToken validates tokens and returns the claim passed into the creator.
func (verifier TestTokenVerifier) ValidateToken(idToken string) (claim *Claim, err error) {
	return verifier.claimToReturn, verifier.errorToReturn
}
