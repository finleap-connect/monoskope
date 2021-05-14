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
	privateKeyFileName string
}

// NewSigner creates a thin wrapper around Square's
// go-jose library to issue JWT.
func NewSigner(privateKeyFilename string) JWTSigner {
	return &jwtSigner{
		privateKeyFileName: privateKeyFilename,
	}
}

// createSigner loads the private key and a returns a new jose.Signer
func (signer *jwtSigner) createSigner() (jose.Signer, error) {
	// Read private key from file
	privKeyBytes, err := ioutil.ReadFile(signer.privateKeyFileName)
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

	return rsaSigner, nil
}

// GenerateSignedToken generates a signed JWT containing the given claims
func (signer *jwtSigner) GenerateSignedToken(claims interface{}) (string, error) {
	joseSigner, err := signer.createSigner()
	if err != nil {
		return "", err
	}

	builder := jwt.Signed(joseSigner).Claims(claims)
	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		return "", err
	}
	return rawJWT, nil
}
