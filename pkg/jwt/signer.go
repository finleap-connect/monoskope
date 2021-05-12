package jwt

import (
	"io/ioutil"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	SignatureAlgorithm = jose.RS512
)

// JWTSigner is an interface for JWT signers
type JWTSigner interface {
	// GenerateSignedToken generates a signed JWT containing the given claims
	GenerateSignedToken(interface{}) (string, error)
}

type jwtSigner struct {
	jose.Signer
}

// NewSigner creates a thin wrapper around Square's
// go-jose library to issue JWT.
func NewSigner(privateKeyFilename string) (JWTSigner, error) {
	// Read private key from file
	privKeyBytes, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		return nil, err
	}

	// Decode
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

	return &jwtSigner{
		Signer: rsaSigner,
	}, nil
}

// GenerateSignedToken generates a signed JWT containing the given claims
func (signer *jwtSigner) GenerateSignedToken(claims interface{}) (string, error) {
	builder := jwt.Signed(signer.Signer).Claims(claims)
	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		return "", err
	}
	return rawJWT, nil
}
