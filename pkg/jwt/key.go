package jwt

import (
	"crypto/x509"
	"encoding/pem"
)

// LoadPublicKey loads a public key from PEM/DER-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {

	block, rest := pem.Decode(data)
	if block != nil {
		data = block.Bytes
	} else if rest != nil {
		data = rest
	}

	pub, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

// LoadPrivateKey loads a private key from PEM/DER-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	block, rest := pem.Decode(data)
	if block != nil {
		data = block.Bytes
	} else if rest != nil {
		data = rest
	}

	if priv, err := x509.ParsePKCS1PrivateKey(data); err != nil {
		return nil, err
	} else {
		return priv, nil
	}
}
