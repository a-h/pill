package tokenverifier

// A TokenVerifier provides a way of validating OAuth ID tokens.
type TokenVerifier interface {
	ValidateToken(idToken string) (claim *Claim, err error)
}
