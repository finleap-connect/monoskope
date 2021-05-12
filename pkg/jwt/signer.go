package jwt

import (
	"io/ioutil"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	SignatureAlgorithm = jose.RS512
)

// NewSigner creates a thin wrapper around Square's
// go-jose library to issue JWT.
func NewSigner(privateKeyFilename string) (jose.Signer, error) {
	privKeyBytes, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		return nil, err
	}

	privKey, err := LoadPrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}

	// create Square.jose signing key
	key := jose.SigningKey{Algorithm: SignatureAlgorithm, Key: privKey}
	// create a Square.jose RSA signer, used to sign the JWT
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")

	rsaSigner, err := jose.NewSigner(key, &signerOpts)
	if err != nil {
		return nil, err
	}
	return rsaSigner, nil
}

// Sign an authentication token and return the serialized JWS
func Sign(signer jose.Signer, claims *jwt.Claims) (string, error) {
	builder := jwt.Signed(signer).Claims(claims)
	signed, err := builder.CompactSerialize()
	if err != nil {
		return "", err
	}
	return signed, nil
}
