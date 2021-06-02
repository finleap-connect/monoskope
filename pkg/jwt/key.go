package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"gopkg.in/square/go-jose.v2"
)

func loadJSONWebKey(json []byte, pub bool) (*jose.JSONWebKey, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON(json)
	if err != nil {
		return nil, err
	}
	if !jwk.Valid() {
		return nil, errors.New("invalid JWK")
	}
	if jwk.IsPublic() != pub {
		return nil, errors.New("priv/pub JWK mismatch")
	}
	return &jwk, nil
}

func convertToJSONWebKey(key interface{}, kid string, pub bool) (*jose.JSONWebKey, error) {
	jwk := jose.JSONWebKey{
		KeyID:     kid,
		Key:       key,
		Algorithm: string(SignatureAlgorithm),
		Use:       "sig",
	}
	if !jwk.Valid() {
		return nil, errors.New("invalid JWK")
	}
	if jwk.IsPublic() != pub {
		return nil, errors.New("priv/pub JWK mismatch")
	}
	return &jwk, nil
}

// LoadPublicKey loads a public key from PEM/DER/JWK-encoded data.
func LoadPublicKey(data []byte) (*jose.JSONWebKey, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	keyID := util.HashBytes(input)

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return convertToJSONWebKey(pub, keyID, true)
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return convertToJSONWebKey(cert, keyID, true)
	}

	jwk, err2 := loadJSONWebKey(data, true)
	if err2 == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s' and '%s'", err0, err1, err2)
}

// LoadPrivateKey loads a private key from PEM/DER/JWK-encoded data.
func LoadPrivateKey(data []byte) (*jose.JSONWebKey, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var err0, err1, err2 error
	if priv, err0 := x509.ParsePKCS1PrivateKey(input); err0 == nil {
		pubKeyPem, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
		if err != nil {
			return nil, err
		}
		return convertToJSONWebKey(priv, util.HashBytes(pubKeyPem), false)
	}

	if priv, err1 := x509.ParsePKCS8PrivateKey(input); err1 == nil {
		switch privTyped := priv.(type) {
		case *rsa.PrivateKey:
			pubKeyPem, err := x509.MarshalPKIXPublicKey(&privTyped.PublicKey)
			if err != nil {
				return nil, err
			}
			return convertToJSONWebKey(priv, util.HashBytes(pubKeyPem), false)
		case *ecdsa.PrivateKey:
			pubKeyPem, err := x509.MarshalPKIXPublicKey(&privTyped.PublicKey)
			if err != nil {
				return nil, err
			}
			return convertToJSONWebKey(priv, util.HashBytes(pubKeyPem), false)
		}
	}

	if priv, err2 := x509.ParseECPrivateKey(input); err2 == nil {
		pubKeyPem, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
		if err != nil {
			return nil, err
		}
		return convertToJSONWebKey(priv, util.HashBytes(pubKeyPem), false)
	}

	jwk, err3 := loadJSONWebKey(input, false)
	if err3 == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s', '%s' and '%s'", err0, err1, err2, err3)
}
