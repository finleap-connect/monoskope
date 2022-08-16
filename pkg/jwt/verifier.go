// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"errors"
	"os"
	"sync"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// JWTVerifier verifies a JWT and parses claims
type JWTVerifier interface {
	Verify(string, interface{}) error
	JWKS() *jose.JSONWebKeySet
	Close()
}

type jwtVerifier struct {
	log        logger.Logger
	jsonWebKey *jose.JSONWebKey
	watcher    *fsnotify.Watcher
	mutex      sync.RWMutex
	watching   chan struct{}
}

// NewVerifier creates a new verifier for raw JWT
func NewVerifier(publicKeyFilename string) (JWTVerifier, error) {
	v := &jwtVerifier{
		log: logger.WithName("jwt-verifier"),
	}

	v.log.Info("Loading public key...", "publicKeyFilename", publicKeyFilename)
	err := v.rotatePublicKey(publicKeyFilename)
	if err != nil {
		return nil, err
	}

	v.log.Info("Setting up watcher...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(publicKeyFilename)
	if err != nil {
		return nil, err
	}
	v.watcher = watcher
	v.watching = make(chan struct{})
	go v.loadPublicKeyOnFileChange()

	return v, nil
}

func (v *jwtVerifier) loadPublicKeyOnFileChange() {
	defer func() {
		defer v.watcher.Close()
		v.log.Info("Watcher closed.")
	}()
	for {
		select {
		case <-v.watching:
			return
		case event := <-v.watcher.Events:
			v.log.Info("Public key has been changed. Updating...")
			err := v.rotatePublicKey(event.Name)
			if err != nil {
				v.log.Error(err, "Error rotating public key.")
			}
			v.log.Info("Public key has been updated.", "KeyID", v.jsonWebKey.KeyID)
		case err := <-v.watcher.Errors:
			v.log.Error(err, "Error from watcher.")
		}
	}
}

// loadPublicKey loads the public key
func (v *jwtVerifier) rotatePublicKey(filename string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	pubKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	pubKey, err := LoadPublicKey(pubKeyBytes)
	if err != nil {
		return err
	}

	v.jsonWebKey = pubKey

	return nil
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims interface{}) error {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	// parse the raw jwt
	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	kid := parsedJWT.Headers[0].KeyID
	if v.jsonWebKey.KeyID == kid {
		if err := parsedJWT.Claims(v.jsonWebKey, claims); err == nil {
			v.log.Info("Successfully verified claims.", "KeyID", v.jsonWebKey.KeyID)
			return nil
		} else {
			v.log.Info("Failed to verify claims.", "KeyID", v.jsonWebKey.KeyID, "error", err.Error())
		}
	}

	// none of the known keys could verify claims
	return errors.New("failed to verify claims")
}

func (v *jwtVerifier) JWKS() *jose.JSONWebKeySet {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	jwks := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0),
	}
	jwks.Keys = append(jwks.Keys, *v.jsonWebKey)
	return jwks
}

// Close closes file watcher
func (v *jwtVerifier) Close() {
	close(v.watching)
}
