package tokenverifier

// A Claim represents a subset of fields available in a JWT.
// See (https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32)
type Claim struct {
	// The issuer, should be "https://accounts.google.com" or "accounts.google.com"
	Issuer string `json:"iss"`
	// The expiry, e.g. "1433981953". Should not be in the past.
	Expiry        string `json:"exp"`
	Email         string `json:"email"` // e.g. "testuser@gmail.com",
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`        // e.g. "Test User",
	Picture       string `json:"picture"`     // e.g. "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
	GivenName     string `json:"given_name"`  // e.g. "Test"
	FamilyName    string `json:"family_name"` // e.g. "User"
}
