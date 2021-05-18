package jwt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Verifier verifies a JWT and parses claims
type Verifier interface {
	Verify(string, interface{}) error
	Close() error
}

type key struct {
	publicKey interface{}
	expiry    time.Time
}

type jwtVerifier struct {
	log           logger.Logger
	keyExpiration time.Duration
	publicKey     []key
	watcher       *fsnotify.Watcher
	mutex         sync.RWMutex
}

// NewVerifier creates a new verifier for raw JWT
func NewVerifier(publicKeyFilename string, keyExpiration time.Duration) (Verifier, error) {
	v := &jwtVerifier{
		log:           logger.WithName("jwt-verifier"),
		publicKey:     make([]key, 0),
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
	go v.loadPublicKeyOnFileChange()

	err = watcher.Add(publicKeyFilename)
	if err != nil {
		return nil, err
	}
	v.watcher = watcher

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
				v.log.Info("Public key has been updated.", "KeyCount", len(v.publicKey))
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

	v.publicKey = append(v.publicKey, key{publicKey: pubKey, expiry: time.Now().UTC().Add(v.keyExpiration)})
	v.removeExpiredKeys()

	return nil
}

// removeExpiredKeys removes expired public keys from cache
func (v *jwtVerifier) removeExpiredKeys() {
	var validKeys []key
	for _, k := range v.publicKey {
		// check if expired and give a little extra time to verify
		if k.expiry.Before(time.Now().UTC().Add(1 * time.Minute)) {
			v.log.Info("Public key expired. Removing from list.", "Expiry", k.expiry, "KeyCount", len(v.publicKey))
		} else {
			validKeys = append(validKeys, k)
		}
	}

	if len(validKeys) == 0 {
		err := fmt.Errorf("not a single valid public key available for verifying claims")
		v.log.Error(err, "no public keys available")
		panic(err)
	}

	v.publicKey = validKeys
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

	// iterate known public keys to validate signature
	for _, pubKey := range v.publicKey {
		if err := parsedJWT.Claims(pubKey.publicKey, claims); err == nil {
			v.log.Info("Successfully verified claims with one of the known public keys.", "Expiry", pubKey.expiry, "KeyCount", len(v.publicKey))
			return nil
		} else {
			v.log.Info("Failed to verify claims with one of the known public keys.", "Expiry", pubKey.expiry, "KeyCount", len(v.publicKey), "error", err.Error())
		}
	}

	// none of the known keys could verify claims
	return errors.New("failed to very claims")
}

// Close closes file watcher
func (v *jwtVerifier) Close() error {
	return v.watcher.Close()
}
