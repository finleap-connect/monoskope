package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// LoadPublicKey loads a public key from PEM/DER-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

// LoadPrivateKey loads a private key from PEM/DER-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var errors []error
	if priv, err := x509.ParsePKCS1PrivateKey(input); err == nil {
		return priv, nil
	} else {
		errors = append(errors, err)
	}

	if priv, err := x509.ParsePKCS8PrivateKey(input); err == nil {
		return priv, nil
	} else {
		errors = append(errors, err)
	}

	if priv, err := x509.ParseECPrivateKey(input); err == nil {
		return priv, nil
	} else {
		errors = append(errors, err)
	}

	return nil, fmt.Errorf("failed to parse private key: %v", errors)
}
