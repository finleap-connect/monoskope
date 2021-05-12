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
	publicKey interface{}
}

// NewVerifier creates a new verifier for raw JWTs
func NewVerifier(publicKeyFilename string) (Verifier, error) {
	pubKeyBytes, err := ioutil.ReadFile(publicKeyFilename)
	if err != nil {
		return nil, err
	}

	pubKey, err := LoadPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return &jwtVerifier{
		publicKey: pubKey,
	}, nil
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims interface{}) error {
	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	if err := parsedJWT.Claims(v.publicKey, claims); err != nil {
		return err
	}

	return nil
}
