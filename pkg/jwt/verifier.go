package jwt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// JWTVerifier verifies a JWT and parses claims
type JWTVerifier interface {
	Verify(string, interface{}) error
	JWKS() *jose.JSONWebKeySet
	KeyExpiration() time.Duration
	Close() error
}

type jwkWithExpiry struct {
	jwk    *jose.JSONWebKey
	expiry time.Time
}

type jwtVerifier struct {
	log           logger.Logger
	keyExpiration time.Duration
	jsonWebKeys   []jwkWithExpiry
	watcher       *fsnotify.Watcher
	mutex         sync.RWMutex
}

// NewVerifier creates a new verifier for raw JWT
func NewVerifier(publicKeyFilename string, keyExpiration time.Duration) (JWTVerifier, error) {
	v := &jwtVerifier{
		log:           logger.WithName("jwt-verifier"),
		jsonWebKeys:   make([]jwkWithExpiry, 0),
		keyExpiration: keyExpiration,
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

	go v.loadPublicKeyOnFileChange()

	return v, nil
}

func (v *jwtVerifier) loadPublicKeyOnFileChange() {
	for {
		select {
		case event, ok := <-v.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				v.log.Info("Public key has been changed. Updating...")
				err := v.rotatePublicKey(event.Name)
				if err != nil {
					v.log.Error(err, "Error rotating public key.")
				}
				v.log.Info("Public key has been updated.", "KeyCount", len(v.jsonWebKeys))
			}
		case err, ok := <-v.watcher.Errors:
			if !ok {
				return
			}
			v.log.Error(err, "Error from watcher.")
		}
	}
}

// loadPublicKey loads the public key
func (v *jwtVerifier) rotatePublicKey(filename string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	pubKeyBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	pubKey, err := LoadPublicKey(pubKeyBytes)
	if err != nil {
		return err
	}

	v.jsonWebKeys = append(v.jsonWebKeys, jwkWithExpiry{jwk: pubKey, expiry: time.Now().UTC().Add(v.keyExpiration)})
	v.removeExpiredKeys()

	return nil
}

// removeExpiredKeys removes expired public keys from cache
func (v *jwtVerifier) removeExpiredKeys() {
	var validKeys []jwkWithExpiry
	for _, k := range v.jsonWebKeys {
		// check if expired and give a little extra time to verify
		if k.expiry.Before(time.Now().UTC().Add(1 * time.Minute)) {
			v.log.Info("Public key expired. Removing from list.", "Expiry", k.expiry, "KeyCount", len(v.jsonWebKeys))
		} else {
			validKeys = append(validKeys, k)
		}
	}

	if len(validKeys) == 0 {
		err := fmt.Errorf("not a single valid public key available for verifying claims")
		v.log.Error(err, "no public keys available")
		panic(err)
	}

	v.jsonWebKeys = validKeys
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims interface{}) error {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	// remove outdated keys
	v.removeExpiredKeys()

	// parse the raw jwt
	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	kid := parsedJWT.Headers[0].KeyID
	for _, key := range v.jsonWebKeys {
		if key.jwk.KeyID == kid {
			if err := parsedJWT.Claims(key.jwk, claims); err == nil {
				v.log.Info("Successfully verified claims.", "Expiry", key.expiry, "KeyCount", len(v.jsonWebKeys), "KeyID", key.jwk.KeyID)
				return nil
			} else {
				v.log.Info("Failed to verify claims.", "Expiry", key.expiry, "KeyCount", len(v.jsonWebKeys), "KeyID", key.jwk.KeyID, "error", err.Error())
			}
		}
	}

	// none of the known keys could verify claims
	return errors.New("failed to very claims")
}

func (v *jwtVerifier) JWKS() *jose.JSONWebKeySet {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	jwks := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0),
	}
	for _, k := range v.jsonWebKeys {
		jwks.Keys = append(jwks.Keys, *k.jwk)
	}
	return jwks
}

func (v *jwtVerifier) KeyExpiration() time.Duration {
	return v.keyExpiration
}

// Close closes file watcher
func (v *jwtVerifier) Close() error {
	return v.watcher.Close()
}
