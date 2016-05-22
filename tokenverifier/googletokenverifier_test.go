package tokenverifier

import (
	"strconv"
	"testing"
	"time"
)

func TestThatAClaimCanBeUnmarshalled(t *testing.T) {
	claimJSON := `{
 "iss": "https://accounts.google.com",
 "sub": "110169484474386276334",
 "azp": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
 "aud": "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
 "iat": "1433978353",
 "exp": "1433981953",
 "email": "testuser@gmail.com",
 "email_verified": "true",
 "name" : "Test User",
 "picture": "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
 "given_name": "Test",
 "family_name": "User",
 "locale": "en"
}`

	claim, err := NewClaim([]byte(claimJSON))

	if err != nil {
		t.Fatal("Received an error during unmarshalling. ", err)
	}

	validate("issuer", "https://accounts.google.com", claim.Issuer, t)
	validate("exp", "1433981953", claim.Expiry, t)
	validate("email", "testuser@gmail.com", claim.Email, t)
	validate("email_verified", "true", claim.EmailVerified, t)
	validate("name", "Test User", claim.Name, t)
	validate("picture", "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg", claim.Picture, t)
	validate("given_name", "Test", claim.GivenName, t)
	validate("family_name", "User", claim.FamilyName, t)
}

func validate(field string, expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("For field %s, expected \"%s\", actual \"%s\"", field, expected, actual)
	}
}

func TestThatAValidClaimIsValidatedSuccessfully(t *testing.T) {
	secondsSince1970 := time.Now().Add(time.Hour).Unix()
	expiry := strconv.Itoa(int(secondsSince1970))

	claim := &Claim{
		Issuer:        "accounts.google.com",
		Expiry:        expiry,
		EmailVerified: "true",
		Email:         "a-h@github.com",
	}

	gtv := &GoogleTokenVerifier{}
	ok, err := gtv.IsClaimValid(claim)

	if !ok {
		t.Error("The claim had all of the required properties set correctly. ", err)
		t.Error(claim)
	}
}

func TestThatExpiredClaimFailValidation(t *testing.T) {
	secondsSince1970 := time.Now().Add(-time.Hour).Unix()
	expiry := strconv.Itoa(int(secondsSince1970))

	claim := &Claim{
		Issuer:        "accounts.google.com",
		Expiry:        expiry,
		EmailVerified: "true",
		Email:         "a-h@github.com",
	}

	gtv := &GoogleTokenVerifier{}
	ok, err := gtv.IsClaimValid(claim)

	if ok {
		t.Error("The claim should not have been passed, its expiry is in the past. ", err)
		t.Error(claim)
	}
}

func TestThatClaimsWithInvalidIssuersFailValidation(t *testing.T) {
	secondsSince1970 := time.Now().Add(time.Hour).Unix()
	expiry := strconv.Itoa(int(secondsSince1970))

	claim := &Claim{
		Issuer:        "https://another.provider.com",
		Expiry:        expiry,
		EmailVerified: "true",
		Email:         "a-h@github.com",
	}

	gtv := &GoogleTokenVerifier{}
	ok, err := gtv.IsClaimValid(claim)

	if ok {
		t.Error("The claim should not have been passed, its issuer is invalid. ", err)
		t.Error(claim)
	}
}

func TestThatClaimsWithoutVerifiedEmailAddressesFailValidation(t *testing.T) {
	inTheFuture := time.Now().Add(time.Hour).Unix()
	expiry := strconv.Itoa(int(inTheFuture))

	claim := &Claim{
		Issuer:        "https://accounts.google.com",
		Expiry:        expiry,
		EmailVerified: "false",
		Email:         "a-h@github.com",
	}

	gtv := &GoogleTokenVerifier{}
	ok, err := gtv.IsClaimValid(claim)

	if ok {
		t.Error("The claim should not have been passed, the email has not been validated. ", err)
		t.Error(claim)
	}
}

func TestThatClaimsWithMissingEmailAddressesFailValidation(t *testing.T) {
	secondsSince1970 := time.Now().Add(time.Hour).Unix()
	expiry := strconv.Itoa(int(secondsSince1970))

	claim := &Claim{
		Issuer:        "https://accounts.google.com",
		Expiry:        expiry,
		EmailVerified: "true",
		Email:         "",
	}

	gtv := &GoogleTokenVerifier{}
	ok, err := gtv.IsClaimValid(claim)

	if ok {
		t.Error("The claim should not have been passed, the email was not set. ", err)
		t.Error(claim)
	}
}
