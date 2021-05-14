package jwt

import (
	"io/ioutil"

	"gopkg.in/square/go-jose.v2/jwt"
)

// Verifier verifies a JWT and parses claims
type Verifier interface {
	Verify(string, interface{}) error
}

type jwtVerifier struct {
	publicKeyFilename string
}

// NewVerifier creates a new verifier for raw JWT
func NewVerifier(publicKeyFilename string) Verifier {
	return &jwtVerifier{
		publicKeyFilename: publicKeyFilename,
	}
}

// loadPublicKey loads the public key
func (v *jwtVerifier) loadPublicKey() (interface{}, error) {
	pubKeyBytes, err := ioutil.ReadFile(v.publicKeyFilename)
	if err != nil {
		return nil, err
	}

	pubKey, err := LoadPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims interface{}) error {
	pubKey, err := v.loadPublicKey()
	if err != nil {
		return err
	}

	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	if err := parsedJWT.Claims(pubKey, claims); err != nil {
		return err
	}

	return nil
}
