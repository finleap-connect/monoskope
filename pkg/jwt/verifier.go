package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"gopkg.in/square/go-jose.v2/jwt"
)

// Verifier verifies a JWT and parses claims
type Verifier interface {
	Verify(string, *interface{}) error
}

type jwtVerifier struct {
	publicKey *rsa.PublicKey
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

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Expected the public key to use ECDSA, but got a key of type %T", pubKey)
	}

	return &jwtVerifier{
		publicKey: rsaPublicKey,
	}, nil
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims *interface{}) error {
	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	if err := parsedJWT.Claims(&v.publicKey, &claims); err != nil {
		return err
	}

	return nil
}
